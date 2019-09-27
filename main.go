/*
MCSS2 - Minecraft Server Save 2
Author: Benjamin Miles
Date: 9.27.2019
Notes: Program to backup select region files and player data to archive.
*/

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dootbin/MCSS/config"
	"github.com/dootbin/MCSS/copy"
	"github.com/dootbin/MCSS/messenger"
	"github.com/jinzhu/now"
	"github.com/mholt/archiver"
)

//LogReport Report that will be sent out.
var LogReport string

// Copy region files
func copyRegions(tmp string, targetDirectory string, saveDiameter int) {

	//16blocks in a chunk 32 chunks in a region bitshift.
	regions := saveDiameter >> 9
	counter := 0
	regionVerifyTarget := (regions + (regions + 2)) * (regions + (regions + 2))

	for i := ((regions + 1) * -1); i <= regions; i++ {

		for a := ((regions + 1) * -1); a <= regions; a++ {

			filename := "r." + strconv.Itoa(i) + "." + strconv.Itoa(a) + ".mca"

			targetFile := targetDirectory + "/" + filename
			saveTarget := tmp + "/" + filename

			err := copy.Copy(targetFile, saveTarget)

			if err != nil {

				LogReport += filename + " FAILED TO COPY\n"

			} else {
				counter++
			}
		}
	}

	LogReport += "Number of regions = " + strconv.Itoa(counter) + "\n"
	LogReport += "Target of regions = " + strconv.Itoa(regionVerifyTarget) + "\n"

	if counter == regionVerifyTarget {
		LogReport += "You have the correct number of region files\n"
	} else {

		LogReport += "You have missed the target number of regions\n"
	}

}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

//archiveDelete crafts and execute curl command to delete file from server
func archiveDelete(path string, filename string) string  {

	ftd := fmt.Sprintf("%s/%s", path, filename)

	archiveExists, err := exists(ftd)

	if archiveExists {
		err = os.Remove(ftd)
		LogReport += filename + " DELETED\n"
	} else {

		LogReport += "Failed to DELETE " + filename + "\n"

	}

	if err != nil {
		return err.Error()
	}
	return ftd
}


func main() {

	//Read Config
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
	}

	saveDiameter, err := strconv.Atoi(config.SaveDiameter)

	//var saveDirectory = config.SaveDir
	var serverRootDirectory = config.ServerRootDirectory
	var serverName = config.ServerName
	var saveName = config.SaveName
	var worldName = config.WorldName
	var webHookURL = config.WebHookURL
	var saveLocation = config.SaveDir
	currentTime := time.Now()
	d := currentTime.Day()
	m := int(currentTime.Month())
	y := currentTime.Year()
	//End of Month
	eom := int(now.EndOfMonth().Day())
	lastMonth := m - 1

	tmp := fmt.Sprintf("/%s/tmp", serverRootDirectory)
	isSaveFolderVerified, _ := exists(saveLocation)
	isTmpFolderVerified, _ := exists(tmp)

	//Verify that there is a clean directory to work in.
	if isTmpFolderVerified {

		os.RemoveAll(tmp)
		os.MkdirAll(tmp, os.ModePerm)

	} else {

		os.MkdirAll(tmp, os.ModePerm)

	}

	if !isSaveFolderVerified {

		os.MkdirAll(saveLocation, os.ModePerm)

	}

	//Save ender = saveDiameter/8
	enderDiameter := saveDiameter
	enderSource := fmt.Sprintf("/%s/%s/%s_the_end/DIM1/region", serverRootDirectory, serverName, worldName)

	//copy ender regions to tmp folder.
	enderDest := fmt.Sprintf("/%s/%s_the_end/DIM1/region", tmp, worldName)
	LogReport += "Attempted to Copy the end\n"
	copyRegions(enderDest, enderSource, enderDiameter)

	//Save Nether = SaveDiameter/8
	netherDiameter := saveDiameter
	netherSource := fmt.Sprintf("/%s/%s/%s_nether/DIM-1/region", serverRootDirectory, serverName, worldName)

	//copy nether regions to tmp folder.
	netherDest := fmt.Sprintf("/%s/%s_nether/DIM-1/region", tmp, worldName)
	LogReport += "Attempted to Copy Nether\n"
	copyRegions(netherDest, netherSource, netherDiameter)

	//copy overworld regions to tmp folder.
	overworldSource := fmt.Sprintf("/%s/%s/%s/region", serverRootDirectory, serverName, worldName)
	overworldDest := fmt.Sprintf("/%s/%s/region", tmp, worldName)
	LogReport += "Attempted to Copy OverWorld\n"
	copyRegions(overworldDest, overworldSource, saveDiameter)

	//Copy player data over to tmp folder
	playerDataLocation := fmt.Sprintf("/%s/%s/%s/playerdata", serverRootDirectory, serverName, worldName)
	playerDataTmpLocation := fmt.Sprintf("/%s/%s/playerdata", tmp, worldName)
	LogReport += "Attempted to Copy PlayerData\n"
	copy.Copy(playerDataLocation, playerDataTmpLocation)

	//Create backup name
	mString := strconv.Itoa(m)
	dString := strconv.Itoa(d)
	yString := strconv.Itoa(y)
	backupName := fmt.Sprintf("%s.%s.%s.%s.tar.gz", saveName, mString, dString, yString)

	//Compress backup
	backupDest := fmt.Sprintf("%s/%s", saveLocation, backupName)
	err = archiver.Archive([]string{tmp}, backupDest)

	//Delete Old Backup. Keep one month's worth of backups and additionally a years worth of monthly backups (12)
	if lastMonth == 0 {
		lastMonth = 12
	}

	var fileToDelete string
	if d == eom {

		deleteCounter := 31 - d
		for i := 0; i <= deleteCounter; i++ {

			fileToDelete = saveName + strconv.Itoa(lastMonth) + "." + strconv.Itoa(d+i) + "." + strconv.Itoa(y) + ".tar.gz"
			_ = archiveDelete(saveLocation, fileToDelete)

		}

	} else {

		if d < 1 {

			fileToDelete = saveName + strconv.Itoa(lastMonth) + "." + strconv.Itoa(d) + "." + strconv.Itoa(y) + ".tar.gz"
			_ = archiveDelete(saveLocation, fileToDelete)

		} else {

			fileToDelete = saveName + strconv.Itoa(m) + "." + strconv.Itoa(d) + "." + strconv.Itoa(y-1) + ".tar.gz"
			_ = archiveDelete(saveLocation, fileToDelete)

		}



	}

	os.RemoveAll(tmp)
	LogReport += "Finished Save"
	messenger.DiscordMessage(LogReport, webHookURL)

}