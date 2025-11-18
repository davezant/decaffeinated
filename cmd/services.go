/*
Copyright © 2025 Davezant <dsndeividdsn1@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/davezant/decafein/src/server/webserver"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

const initialSecret = "ineedtochangethissecret"

// -------------------------------------------------------------------
// Programa que o serviço executa
// -------------------------------------------------------------------

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	log.Println("Servidor Decafein iniciado como serviço...")

	go func() {
		// O servidor RODA aqui e bloqueia AQUI dentro,
		// mas não bloqueia o serviço
		webserver.OpenServer(initialSecret, true)
	}()

	// Mantém a função viva sem bloquear
	// (Windows exige que o serviço não finalize)
	select {}
}

func (p *program) Stop(s service.Service) error {
	log.Println("Servidor Decafein parado.")
	return nil
}

// -------------------------------------------------------------------
// Utilidade para criar serviço
// -------------------------------------------------------------------

func newService() service.Service {
	svcConfig := &service.Config{
		Name:        "DecafeinServer",
		DisplayName: "Decafein Server",
		Description: "Servidor de controle de atividades e processos Decafein",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatalf("Erro ao criar serviço: %v", err)
	}
	return s
}

// -------------------------------------------------------------------
// Comandos do serviço
// -------------------------------------------------------------------

var serviceCmd = &cobra.Command{
	Use:   "services",
	Short: "Controla o serviço Decafein",
	Long:  `Controla o serviço do Decafein usando kardianos/service.`,
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Instala o serviço system-wide",
	Run: func(cmd *cobra.Command, args []string) {
		s := newService()
		if err := s.Install(); err != nil {
			log.Fatalf("Erro ao instalar: %v", err)
		}
		fmt.Println("Serviço instalado com sucesso.")
	},
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove o serviço",
	Run: func(cmd *cobra.Command, args []string) {
		s := newService()
		if err := s.Uninstall(); err != nil {
			log.Fatalf("Erro ao desinstalar: %v", err)
		}
		fmt.Println("Serviço removido com sucesso.")
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o serviço",
	Run: func(cmd *cobra.Command, args []string) {
		s := newService()
		if err := s.Start(); err != nil {
			log.Fatalf("Erro ao iniciar: %v", err)
		}
		fmt.Println("Serviço iniciado.")
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Para o serviço",
	Run: func(cmd *cobra.Command, args []string) {
		s := newService()
		if err := s.Stop(); err != nil {
			log.Fatalf("Erro ao parar: %v", err)
		}
		fmt.Println("Serviço parado.")
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Reinicia o serviço",
	Run: func(cmd *cobra.Command, args []string) {
		s := newService()
		if err := s.Restart(); err != nil {
			log.Fatalf("Erro ao reiniciar: %v", err)
		}
		fmt.Println("Serviço reiniciado.")
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Executa o servidor em modo foreground (sem serviço)",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Servidor Decafein iniciado no terminal...")
		webserver.OpenServer(initialSecret, false)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Mostra o status do serviço",
	Run: func(cmd *cobra.Command, args []string) {
		s := newService()
		status, err := s.Status()
		if err != nil {
			fmt.Println("Não foi possível obter status:", err)
			return
		}

		switch status {
		case service.StatusRunning:
			fmt.Println("Status: RODANDO")
		case service.StatusStopped:
			fmt.Println("Status: PARADO")
		default:
			fmt.Println("Status: DESCONHECIDO")
		}
	},
}

// -------------------------------------------------------------------
// Comando EXTRA: open
// Inicia o serviço e abre o navegador
// -------------------------------------------------------------------

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Inicia o serviço e abre a interface web",
	Run: func(cmd *cobra.Command, args []string) {
		s := newService()

		// instala se não existir
		if err := s.Install(); err != nil && err.Error() != "service already exists" {
			log.Fatalf("Erro ao instalar serviço: %v", err)
		}

		// inicia
		if err := s.Start(); err != nil && err.Error() != "the service is already running" {
			log.Fatalf("Erro ao iniciar serviço: %v", err)
		}

		// abre interface via navegador
		fmt.Println("Abrindo interface web...")
		webserver.OpenServer(initialSecret, false)

		fmt.Println("Serviço iniciado e interface carregada.")
	},
}

// -------------------------------------------------------------------
// Inicialização
// -------------------------------------------------------------------

func init() {
	if globalFlags.IsAdmin {

		rootCmd.AddCommand(serviceCmd)

		serviceCmd.AddCommand(installCmd)
		serviceCmd.AddCommand(uninstallCmd)
		serviceCmd.AddCommand(startCmd)
		serviceCmd.AddCommand(stopCmd)
		serviceCmd.AddCommand(restartCmd)
		serviceCmd.AddCommand(statusCmd)
		serviceCmd.AddCommand(runCmd)
		serviceCmd.AddCommand(openCmd)
	}
}
