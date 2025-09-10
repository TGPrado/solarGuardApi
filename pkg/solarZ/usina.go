package solarz

import (
	"fmt"
	"time"

	"github.com/TGPrado/GuardIA/config"
)

func GetUsinaId(id int64, configSolarZ config.SolarZ) (int64, error) {
	for i := 1; i <= 10; i++ {
		usina, err := GetPlants(fmt.Sprintf("%d", id), configSolarZ)
		fmt.Println(usina, err)
		if err != nil {
			return 0, err
		}

		if usina[0].UsinaJaImportadaID == nil {
			time.Sleep(2 * time.Second)
			continue
		}

		return *usina[0].UsinaJaImportadaID, nil
	}

	return 0, fmt.Errorf("error pegando id da usina")
}
