package service

import (
	"context"
	"encoding/hex"
	"errors"
	"strings"

	"video_conference/pkg/logger"

	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/codec/opus"
	"github.com/pion/mediadevices/pkg/prop"

	// Регистрируем драйверы для аудио на Linux
	_ "github.com/pion/mediadevices/pkg/driver/microphone" // Драйвер микрофона для Linux (использует malgo)
)

type AudioCaptureService interface {
	StartCapture(ctx context.Context) (mediadevices.MediaStream, error)
	StopCapture()
}

type audioCaptureService struct {
	log        logger.Logger
	stream     mediadevices.MediaStream
	cancelFunc context.CancelFunc
}

func NewAudioCaptureService(log logger.Logger) AudioCaptureService {
	return &audioCaptureService{
		log: log,
	}
}

func (s *audioCaptureService) StartCapture(ctx context.Context) (mediadevices.MediaStream, error) {
	s.log.Info("Starting audio capture")

	// Настройка Opus кодера для аудио
	opusParams, err := opus.NewParams()
	if err != nil {
		s.log.Error("Failed to create Opus params", "error", err)
		return nil, err
	}
	opusParams.BitRate = 64000 // 64 kbps

	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithAudioEncoders(&opusParams),
	)

	// Получаем список доступных аудио устройств
	devices := mediadevices.EnumerateDevices()
	s.log.Info("Available media devices", "count", len(devices))

	var audioInputDevices []mediadevices.MediaDeviceInfo
	var loopbackDevice *mediadevices.MediaDeviceInfo

	for _, device := range devices {
		if device.Kind == mediadevices.AudioInput {
			audioInputDevices = append(audioInputDevices, device)

			// Декодируем hex label для читаемости
			label := device.Label
			decodedLabel := label
			if strings.HasPrefix(label, "50756c7365") || strings.HasPrefix(label, "616c7361") {
				// Это hex-encoded label, попробуем декодировать
				if decoded, err := hex.DecodeString(label); err == nil {
					decodedLabel = string(decoded)
				}
			}

			s.log.Info("Found audio input device",
				"device_id", device.DeviceID,
				"label", decodedLabel)

			// Ищем monitor устройство (для системного звука)
			// Приоритет: apps.monitor > monitor > output.monitor
			labelLower := strings.ToLower(decodedLabel)
			rawLabelLower := strings.ToLower(label)

			// Проверяем в декодированном и raw label (hex)
			// "monitor" в hex = "6d6f6e69746f72"
			isMonitor := strings.Contains(labelLower, "monitor") || strings.Contains(rawLabelLower, "6d6f6e69746f72")
			isAppsMonitor := strings.Contains(labelLower, "apps.monitor") || strings.Contains(labelLower, "app.monitor") ||
				strings.Contains(rawLabelLower, "617070732e6d6f6e69746f72") || strings.Contains(rawLabelLower, "6170702e6d6f6e69746f72")

			if isMonitor {
				if loopbackDevice == nil || isAppsMonitor {
					loopbackDevice = &device
					if isAppsMonitor {
						s.log.Info("Found apps.monitor device (PREFERRED for system audio)", "label", decodedLabel)
					} else {
						s.log.Info("Found monitor device for system audio", "label", decodedLabel)
					}
				}
			} else if loopbackDevice == nil && (strings.Contains(labelLower, "loopback") || strings.Contains(labelLower, "output")) {
				loopbackDevice = &device
				s.log.Info("Found potential loopback device", "label", decodedLabel)
			}
		}
	}

	s.log.Info("Audio input devices found", "count", len(audioInputDevices))

	// Получаем аудио поток - используем устройство по умолчанию для захвата системного звука
	// Приоритет: loopback/monitor (системный звук) > устройство по умолчанию
	s.log.Info("Requesting audio stream from system (default device for system audio)")

	var mediaStream mediadevices.MediaStream
	var selectedDevice *mediadevices.MediaDeviceInfo

	// Сначала пробуем loopback/monitor устройство (системный звук)
	if loopbackDevice != nil {
		decodedLabel := loopbackDevice.Label
		if decoded, err := hex.DecodeString(loopbackDevice.Label); err == nil {
			decodedLabel = string(decoded)
		}
		s.log.Info("Using loopback/monitor device for system audio",
			"device_id", loopbackDevice.DeviceID,
			"label", decodedLabel)

		mediaStream, err = mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
			Audio: func(constraint *mediadevices.MediaTrackConstraints) {
				constraint.DeviceID = prop.String(loopbackDevice.DeviceID)
				constraint.SampleRate = prop.Int(48000) // 48 kHz
				constraint.ChannelCount = prop.Int(2)   // Стерео
			},
			Codec: codecSelector,
		})

		if err == nil {
			selectedDevice = loopbackDevice
			s.log.Info("Successfully using loopback device for system audio")
		} else {
			s.log.Warn("Failed to use loopback device, trying default", "error", err)
		}
	}

	// Если loopback не сработал, ищем ЛЮБОЕ monitor устройство (не микрофон!)
	// НИКОГДА не используем устройство по умолчанию - оно может быть микрофоном!
	if mediaStream == nil {
		s.log.Info("Loopback device failed, searching for ANY monitor device (system audio)")

		// Ищем любое monitor устройство в списке
		for _, device := range audioInputDevices {
			decodedLabel := device.Label
			if decoded, err := hex.DecodeString(device.Label); err == nil && len(decoded) > 0 {
				decodedLabel = string(decoded)
			}

			labelLower := strings.ToLower(decodedLabel)
			rawLabelLower := strings.ToLower(device.Label)

			// Проверяем, что это monitor (системный звук), а НЕ микрофон
			isMonitor := strings.Contains(labelLower, "monitor") || strings.Contains(rawLabelLower, "6d6f6e69746f72")
			isMicrophone := strings.Contains(labelLower, "mic") || strings.Contains(labelLower, "microphone") ||
				strings.Contains(labelLower, "input") && !strings.Contains(labelLower, "output")

			if isMonitor && !isMicrophone {
				s.log.Info("Trying monitor device for system audio",
					"device_id", device.DeviceID,
					"label", decodedLabel)

				mediaStream, err = mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
					Audio: func(constraint *mediadevices.MediaTrackConstraints) {
						constraint.DeviceID = prop.String(device.DeviceID)
						constraint.SampleRate = prop.Int(48000)
						constraint.ChannelCount = prop.Int(2)
					},
					Codec: codecSelector,
				})

				if err == nil {
					selectedDevice = &device
					s.log.Info("Successfully using monitor device for system audio", "label", decodedLabel)
					break
				} else {
					s.log.Warn("Failed to use monitor device", "error", err, "label", decodedLabel)
				}
			}
		}
	}

	// Если все еще нет потока, это критическая ошибка - не используем микрофон!
	if mediaStream == nil {
		s.log.Error("Failed to get system audio (monitor device). Cannot use microphone!")
		return nil, errors.New("failed to get system audio device - only monitor devices are allowed for screen sharing")
	}

	if err != nil {
		s.log.Error("Failed to get user media for audio", "error", err)
		return nil, err
	}

	if mediaStream == nil {
		s.log.Error("MediaStream is nil after GetUserMedia")
		return nil, errors.New("failed to create audio stream")
	}

	audioTracks := mediaStream.GetAudioTracks()
	deviceLabel := "default"
	if selectedDevice != nil {
		deviceLabel = selectedDevice.Label
	}
	s.log.Info("Audio capture started successfully",
		"total_tracks", len(mediaStream.GetTracks()),
		"audio_tracks", len(audioTracks),
		"device", deviceLabel)

	if len(audioTracks) == 0 {
		s.log.Error("No audio tracks in stream after GetUserMedia!")
		return nil, errors.New("no audio tracks in stream")
	}

	for i, at := range audioTracks {
		s.log.Info("Audio track in stream",
			"index", i,
			"track_id", at.ID(),
			"kind", at.Kind(),
			"stream_id", at.StreamID())
	}

	s.stream = mediaStream
	return mediaStream, nil
}

func (s *audioCaptureService) StopCapture() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
	if s.stream != nil {
		// Закрываем все треки в стриме
		for _, track := range s.stream.GetTracks() {
			track.Close()
		}
	}
}
