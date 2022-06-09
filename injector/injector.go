package injector

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/steromano87/harkonnen/load"
	"github.com/steromano87/harkonnen/shooter"
	"io"
	"net"
	"sync"
	"time"
)

type Injector struct {
	Context context.Context
	Logger  zerolog.Logger

	SetUpScript    shooter.Script
	MainScripts    []shooter.Script
	TearDownScript shooter.Script
	loadProfiles   []load.Profile

	shooters  []shooter.Shooter
	waitGroup sync.WaitGroup

	cancelFunc context.CancelFunc

	settings    Settings
	tcpListener net.Listener
}

func New(ctx context.Context, logWriter io.Writer, settings Settings) *Injector {
	output := new(Injector)
	output.Context, output.cancelFunc = context.WithCancel(ctx)

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	output.Logger = zerolog.New(logWriter).With().Timestamp().Logger()
	output.settings = settings

	return output
}

func (i *Injector) Start() {
	// Start TCP server
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", i.settings.BindAddress, i.settings.Port))
	if err != nil {
		panic(err)
	}
	i.tcpListener = listener

}

func (i *Injector) Stop() {
	err := i.tcpListener.Close()
	if err != nil {
		panic(err)
	}
}

func (i *Injector) AddLoadProfile(profile load.Profile) {
	i.loadProfiles = append(i.loadProfiles, profile)
}

func (i *Injector) ExpectedShooters(elapsed time.Duration) int {
	expected := 0

	for _, ramp := range i.loadProfiles {
		expected += ramp.At(elapsed)
	}

	return expected
}

func (i *Injector) AdjustScheduling(elapsed time.Duration) error {
	expectedShooters := i.ExpectedShooters(elapsed)
	activeShooters := i.ActiveShooters()

	var err error

	if activeShooters < expectedShooters {
		err = i.addShooter()
	}

	if activeShooters > expectedShooters {
		err = i.removeShooter()
	}

	return err
}

func (i *Injector) ActiveShooters() int {
	return len(i.shooters)
}

func (i *Injector) addShooter() error {
	newShooter := i.initShooter()
	i.shooters = append(i.shooters, newShooter)
	newShooter.Start()

	return nil
}

func (i *Injector) removeShooter() error {
	shooterToStop := i.shooters[0]
	i.shooters = append(i.shooters[:0], i.shooters[1:]...)

	shooterToStop.ScheduleShutDown()

	return nil
}

func (i *Injector) initShooter() shooter.Shooter {
	shooterID := uuid.NewString()
	shooterLogger := log.With().Str("ID", shooterID).Logger()

	newShooter := shooter.Shooter{
		Context:        shooter.NewContext(i.Context, shooterLogger, shooterID),
		SetUpScript:    i.SetUpScript,
		MainScripts:    i.MainScripts,
		TearDownScript: i.TearDownScript,
		MaxIterations:  0,
		WaitGroup:      &i.waitGroup,
	}

	return newShooter
}
