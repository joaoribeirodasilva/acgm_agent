package meters

import (
	"errors"
	"time"

	"biqx.com.br/acgm_agent/modules/config"
	"github.com/rs/zerolog/log"
)

type Processes struct {
	StartTime      time.Time
	EndTime        time.Time
	Interval       time.Duration
	Processes      []Process
	running        bool
	stop_requested bool
}

func NewProcesses() *Processes {

	m := new(Processes)
	return m
}

func (m *Processes) Start(conf *config.Config) error {

	m.Interval = time.Duration(conf.Metrics.Processes.CollectInterval) * time.Millisecond
	m.stop_requested = false

	go m.loop()

	return nil
}

func (m *Processes) Stop() error {

	m.stop_requested = true

	for {
		attempts := 0

		if !m.running {
			break
		}

		time.Sleep(500 * time.Millisecond)
		attempts++
		if attempts == 4 {
			errStr := "Process meters timeout exiting"
			log.Error().Str("package", "Meters").Msg(errStr)
			return errors.New(errStr)
		}
	}

	return nil
}

// Loop through interface ??

func (m *Processes) loop() {

	m.running = true

	// procFS, err := procfs.NewDefaultFS()
	// if err != nil {
	// 	log.Error().Str("package", "Metrics").Err(err).Msg("Failed to read processes")
	// }

	// for {

	// 	time.Sleep(1000 * time.Millisecond)

	// 	if m.stop_requested {
	// 		break
	// 	}

	// 	// Host processes
	// 	procs, err := procFS.AllProcs()
	// 	if err != nil {
	// 		log.Error().Str("package", "Metrics").Err(err).Msg("Failed to read processes")
	// 		continue
	// 	}

	// 	fmt.Printf("Total Processes: %d\n", procs.Len())

	// 	cpuInfo, err := procFS.CPUInfo()
	// 	if err != nil {
	// 		log.Error().Str("package", "Metrics").Err(err).Msg("Failed to get host CPU information")
	// 		continue
	// 	}

	// 	for _, cpu := range cpuInfo {
	// 		fmt.Printf("CPU Information: \nProcessor: %d\nVendorID: %s\nCPUFamily: %s\nModel: %s\nModel Name: %s, Stepping: %s\nMicrocode: %s\nCPUMHz: %f\nCacheSize: %s\nPhysicalID: %s\nSiblings: %d\nCoreID: %s\nCPUCores: %d\nAPICID: %s\nInitialAPICID: %s\nFPU: %s\nFPUException: %s\nCPUIDLevel: %d\nWP: %s\nFlags: %s\nBugs: %s\nBogoMips: %f\nCLFlushSize: %d\nCacheAlignment: %d\nAddressSizes: %s\nPowerManagement: %s\n",
	// 			cpu.Processor,                 //uint
	// 			cpu.VendorID,                  //string
	// 			cpu.CPUFamily,                 //string
	// 			cpu.Model,                     //string
	// 			cpu.ModelName,                 //string
	// 			cpu.Stepping,                  //string
	// 			cpu.Microcode,                 //string
	// 			cpu.CPUMHz,                    //float64
	// 			cpu.CacheSize,                 //string
	// 			cpu.PhysicalID,                //string
	// 			cpu.Siblings,                  //uint
	// 			cpu.CoreID,                    //string
	// 			cpu.CPUCores,                  //uint
	// 			cpu.APICID,                    //string
	// 			cpu.InitialAPICID,             //string
	// 			cpu.FPU,                       //string
	// 			cpu.FPUException,              //string
	// 			cpu.CPUIDLevel,                //uint
	// 			cpu.WP,                        //string
	// 			strings.Join(cpu.Flags, ", "), //[]string
	// 			strings.Join(cpu.Bugs, ", "),  //[]string
	// 			cpu.BogoMips,                  //float64
	// 			cpu.CLFlushSize,               //uint
	// 			cpu.CacheAlignment,            //uint
	// 			cpu.AddressSizes,              //string
	// 			cpu.PowerManagement,           //string
	// 		)
	// 	}

	// 	netDev, err := procFS.NetDev()
	// 	if err != nil {
	// 		log.Error().Str("package", "Metrics").Err(err).Msg("Failed to read network devices information")
	// 	}

	// 	for key, val := range netDev {
	// 		fmt.Printf("-------------------------\nKey: %s\nInterface: %s, \nRxBytes: %d, \nRxPackets: %d, \nRxErrors: %d, \nRxDropped: %d, \nRxFIFO: %d, \nRxFrame: %d, \nRxCompressed: %d, \nRxMulticast: %d, \nTxBytes: %d, \nTxPackets: %d, \nTxErrors: %d, \nTxDropped: %d, \nTxFIFO: %d, \nTxCollisions: %d, \nTxCarrier: %d, \nTxCompressed: %d\n",
	// 			key,
	// 			val.Name,
	// 			val.RxBytes,
	// 			val.RxPackets,
	// 			val.RxErrors,
	// 			val.RxDropped,
	// 			val.RxFIFO,
	// 			val.RxFrame,
	// 			val.RxCompressed,
	// 			val.RxMulticast,
	// 			val.TxBytes,
	// 			val.TxPackets,
	// 			val.TxErrors,
	// 			val.TxDropped,
	// 			val.TxFIFO,
	// 			val.TxCollisions,
	// 			val.TxCarrier,
	// 			val.TxCompressed,
	// 		)
	// 	}

	// 	netProtocols, err := procFS.NetProtocols()
	// 	if err != nil {
	// 		log.Error().Str("package", "Metrics").Err(err).Msg("Failed to read network protocols information")
	// 	}
	// 	for key, val := range netProtocols {
	// 		fmt.Printf("-------------------------\nKey: %s\nName: %s\nSize: %d\nSockets: %d\nMemory: %d\nPressure: %d\nMaxHeader: %d\nSlab: %t\nModuleName: %s\n",
	// 			key,
	// 			val.Name,          //string // 0 The name of the protocol
	// 			val.Size,          //uint64 // 1 The size, in bytes, of a given protocol structure. e.g. sizeof(struct tcp_sock) or sizeof(struct unix_sock)
	// 			val.Sockets,       //int64  // 2 Number of sockets in use by this protocol
	// 			val.Memory*4*1024, //int64  // 3 Number of 4KB pages allocated by all sockets of this protocol
	// 			val.Pressure,      //int    // 4 This is either yes, no, or NI (not implemented). For the sake of simplicity we treat NI as not experiencing memory pressure.
	// 			val.MaxHeader,     //uint64 // 5 Protocol specific max header size
	// 			val.Slab,          //bool   // 6 Indicates whether or not memory is allocated from the SLAB
	// 			val.ModuleName,    //string // 7 The name of the module that implemented this protocol or "kernel" if not from a module
	// 		)
	// 	}
	//for _, proc := range procs {

	// Executable path
	// exec, err := proc.Executable()
	// if err != nil {
	// 	log.Error().Str("package", "Metrics").Err(err).Msg("Failed to get process executable")
	// }

	// // Extract executable name
	// if !strings.HasSuffix(exec, "/nginx") {
	// 	continue
	// }

	// fmt.Printf("PID: %d, Process: %s\n", proc.PID, exec)

	// process threads
	// threads, err := procFS.AllThreads(proc.PID)
	// if err != nil {
	// 	log.Error().Str("package", "Metrics").Err(err).Msg("Failed to get process threads")
	// }
	// for _, thread := range threads {
	// 	fmt.Printf("PID: %d, Thread PID: %d\n", proc.PID, thread.PID)
	// }

	// Process resource usage
	// stat, err := proc.Stat()
	// if err != nil {
	// 	log.Error().Str("package", "Metrics").Err(err).Msg("Failed to get process stat")
	// }
	// if stat.Comm != "nginx" {
	// 	continue
	// }

	// fmt.Printf("Process PID: %d, PPID: %d, Name: %s, State: %s, CPU: %f, Threads: %d, Mem Res: %d, Mem Virt: %d\n",
	// 	stat.PID,
	// 	stat.PPID,
	// 	stat.Comm,
	// 	stat.State,
	// 	stat.CPUTime(),
	// 	stat.NumThreads,
	// 	stat.ResidentMemory(),
	// 	stat.VirtualMemory(),
	// )

	// Process network usage
	// netDev, err := proc.NetDev()
	// if err != nil {
	// 	log.Error().Str("package", "Metrics").Err(err).Msg("Failed to get process network statistics")
	// }

	// Add totals

	// for _, val := range netDev {
	// 	fmt.Printf("Interface: %s, \nRxBytes: %d, \nRxPackets: %d, \nRxErrors: %d, \nRxDropped: %d, \nRxFIFO: %d, \nRxFrame: %d, \nRxCompressed: %d, \nRxMulticast: %d, \nTxBytes: %d, \nTxPackets: %d, \nTxErrors: %d, \nTxDropped: %d, \nTxFIFO: %d, \nTxCollisions: %d, \nTxCarrier: %d, \nTxCompressed: %d\n",
	// 		val.Name,
	// 		val.RxBytes,
	// 		val.RxPackets,
	// 		val.RxErrors,
	// 		val.RxDropped,
	// 		val.RxFIFO,
	// 		val.RxFrame,
	// 		val.RxCompressed,
	// 		val.RxMulticast,
	// 		val.TxBytes,
	// 		val.TxPackets,
	// 		val.TxErrors,
	// 		val.TxDropped,
	// 		val.TxFIFO,
	// 		val.TxCollisions,
	// 		val.TxCarrier,
	// 		val.TxCompressed,
	// 	)
	// }

	// Process IO usage
	// io, err := proc.IO()
	// if err != nil {
	// 	log.Error().Str("package", "Metrics").Err(err).Msg("Failed to get process IO statistics")
	// }

	// fmt.Printf("Read bytes: %d, \nRead bytes: %d,\n",
	// 	io.ReadBytes,
	// 	io.WriteBytes,
	// )
	// }

	// List all processes
	// Loop through processes
	// Filter process by path (as in config)
	// Get several process metrics (as in config)
	// Store process metrics in the database
	//}

	m.stop_requested = false
	m.running = false
}
