package main

import (
    "log"

    "megabaseGo/internal/config"
    "megabaseGo/internal/models"
    dbpkg "megabaseGo/internal/database"
    dbseed "megabaseGo/internal/database/seeders"

    "github.com/spf13/cobra"
)

func main() {
    var withSeed bool

    rootCmd := &cobra.Command{
        Use:   "console",
        Short: "Herramientas de consola para MegabaseGo",
    }

    migrateCmd := &cobra.Command{
        Use:   "migrate",
        Short: "Ejecuta migraciones y opcionalmente seeders",
        Run: func(cmd *cobra.Command, args []string) {
            // 1) Carga configuración
            cfg := config.LoadConfig()

            // 2) Inicializa conexión a BD
            db, err := dbpkg.InitDB(cfg)
            if err != nil {
                log.Fatalf("Error iniciando BD: %v", err)
            }
            defer dbpkg.CloseDB()

            // 3) Migraciones automáticas basadas en los modelos
            if err := db.AutoMigrate(models.AllModels...); err != nil {
                log.Fatalf("Error en AutoMigrate: %v", err)
            }
            log.Println("✔ Migraciones completadas")

            // 4) Si se pasa --seed, ejecuta todos los seeders
            if withSeed {
                seeder := &dbseed.DatabaseSeeder{}
                if err := seeder.Run(db); err != nil {
                    log.Fatalf("Error ejecutando seeders: %v", err)
                }
                log.Println("✔ Seeders ejecutados")
            }
        },
    }
    migrateCmd.Flags().BoolVarP(&withSeed, "seed", "s", false, "Ejecutar seeders tras migrar")
    rootCmd.AddCommand(migrateCmd)

    if err := rootCmd.Execute(); err != nil {
        log.Fatal(err)
    }
}