package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gabutlabs/devopin/internal/model"
	service "github.com/gabutlabs/devopin/internal/services"
)

type ThresholdAutomation struct {
	workerService       service.WorkerServiceService
	systemMetricService service.SystemMetricService
	alarmService        service.AlarmService
	settingService      service.SettingService
}

// NewThresholdAutomation creates a new ThresholdAutomation instance.
func NewThresholdAutomation(ws service.WorkerServiceService, sm service.SystemMetricService, als service.AlarmService, ss service.SettingService) *ThresholdAutomation {
	return &ThresholdAutomation{
		workerService:       ws,
		systemMetricService: sm,
		alarmService:        als,
		settingService:      ss,
	}
}

func (ta *ThresholdAutomation) AlarmWorker() {
	ta.metricAlarm()
	ta.alarmWorkerService()
}

func (ta *ThresholdAutomation) metricAlarm() {
	settings, err := ta.settingService.GetSettings()
	if err != nil {
		log.Printf("Error get settings: %v", err)
		return
	}

	metric, err := ta.systemMetricService.GetLastSystemMetricFilter("1 hour")
	if err != nil {
		log.Printf("Error get last system metric: %v", err)
	}
	if metric.AvgCPUUsage > settings.SystemCPUCriticalPercent {
		activeAlarm, err := ta.alarmService.CreateActiveAlarm("CPU_ALARM", "system:cpu", model.StatusFiring, fmt.Sprintf("CPU usage last 1 hour greater than %.2f%%", settings.SystemCPUCriticalPercent))
		if err != nil {
			log.Printf("Error create active alarm cpu: %v", err)
		}
		metaData, err := json.Marshal(activeAlarm)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
		}
		ta.alarmService.CreateAlarmHistory("CPU_ALARM", "system:cpu", model.AlarmStatusFiring, activeAlarm.Message, metaData)
		// do alarm notification to email,telegram or whatsapp
	}
	if (metric.AvgDiskUsage/metric.AvgDiskTotal)*100 > settings.SystemDiskCriticalPercent {
		activeAlarm, err := ta.alarmService.CreateActiveAlarm("DISK_ALARM", "system:disk", model.StatusFiring, fmt.Sprintf("DISK usage greater than %.2f%%", settings.SystemDiskCriticalPercent))
		if err != nil {
			log.Printf("Error create active alarm disk: %v", err)

		}
		metaData, err := json.Marshal(activeAlarm)
		if err != nil {
			log.Println("Error marshalling JSON:", err)
		}
		ta.alarmService.CreateAlarmHistory("DISK_ALARM", "system:disk", model.AlarmStatusFiring, activeAlarm.Message, metaData)
		// do alarm notification to email,telegram or whatsapp

	}
	if metric.AvgMemUsage > settings.SystemMemCriticalPercent {
		activeAlarm, err := ta.alarmService.CreateActiveAlarm("MEMORY_ALARM", "system:memory", model.StatusFiring, fmt.Sprintf("Memory usage last 1 hour greater than %.2f%%", settings.SystemMemCriticalPercent))
		if err != nil {
			log.Printf("Error create active alarm memory: %v", err)
		}
		metaData, err := json.Marshal(activeAlarm)
		if err != nil {
			log.Println("Error marshalling JSON:", err)
		}
		ta.alarmService.CreateAlarmHistory("MEMORY_ALARM", "system:memory", model.AlarmStatusFiring, activeAlarm.Message, metaData)
		// do alarm notification to email,telegram or whatsapp
	}
}

func (ta *ThresholdAutomation) alarmWorkerService() {
	settings, err := ta.settingService.GetSettings()
	if err != nil {
		log.Printf("Error get settings: %v", err)
		return
	}

	filters := map[string]map[string]string{}
	workers, err := ta.workerService.ListWorkerServices(filters)
	if err != nil {
		log.Printf("Error fetch worker list: %v", err)
	}
	heartbeatTimeout := time.Duration(settings.WorkerHeartbeatTimeoutSeconds) * time.Second
	for _, v := range workers {
		if v.LastHeartbeatAt != nil && time.Since(*v.LastHeartbeatAt) > heartbeatTimeout {
			activeAlarm, err := ta.alarmService.CreateActiveAlarm(fmt.Sprintf("WORKER: %s", v.Name), "system:worker", model.StatusFiring, fmt.Sprintf("Worker %s inactive grather than %d second", v.Name, heartbeatTimeout/time.Second))
			if err != nil {
				log.Printf("Error insert alarm worker: %v", err)
			}
			metaData, err := json.Marshal(activeAlarm)
			if err != nil {
				log.Println("Error marshalling JSON:", err)
			}
			ta.alarmService.CreateAlarmHistory(fmt.Sprintf("WORKER: %s", v.Name), "system:worker", model.AlarmStatusFiring, activeAlarm.Message, metaData)
		} else if v.CurrentStatus == model.StatusFailed {
			activeAlarm, err := ta.alarmService.CreateActiveAlarm(fmt.Sprintf("WORKER: %s", v.Name), "system:worker", model.StatusFiring, fmt.Sprintf("Worker %s is in failed state", v.Name))
			if err != nil {
				log.Printf("Error insert alarm worker: %v", err)
			}
			metaData, err := json.Marshal(activeAlarm)
			if err != nil {
				log.Println("Error marshalling JSON:", err)
			}
			ta.alarmService.CreateAlarmHistory(fmt.Sprintf("WORKER: %s", v.Name), "system:worker", model.AlarmStatusFiring, activeAlarm.Message, metaData)
		}
	}
}
