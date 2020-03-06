# IPSD
IPSD(Inter Planet Site Watchdog) is a tool to work with IPSC and create static html site automatically.

## Background

IPFS (Inter Planet File System [IPFS](https://ipfs.io)) is a peer-to-peer hyperlink protocol which is used to publish content. We can publish a web site  on IPFS as we publish a site on http.

But as IPFS is an p2p system, file published on IPFS cannot be changed, if we changed a file and publish to IPFS again, it is a completely new file from the old one.  Changing files of a IPFS file is not encouraged. So generally sites that are built on ASP.NET Java PHP which have a lot of scripts are not the best option when you want to publish a site to IPFS. Static website based on HTML and CSS is the best option.

IPSC is the tool to create static html site that you can publish to IPFS.

IPSD work with IPSC to create site automatically.



## Install

Download the release for your platform from Releases, unzip it.


## Build
If you can not find a release for your platform, build it from source code as follows:

1. Install go

2. Install git
   
       	Download and install
       		https://git-scm.com/download
       	OR
       		sudo apt-get install git	

3. Install mingw(Windows)

4. Install Liteide (https://github.com/visualfc/liteide)


   ​	*Windows/Linux/MacOSX just download and install

   ​	*Raspbian

   ​		Download source (qt4 Linux 64)code and compile as follows:

   ​		

   ```bash
       sudo apt-get update
       sudo apt-get upgrade
       sudo apt-get install git
       git clone https://github.com/visualfc/liteide.git
       sudo apt-get install qt4-dev-tools libqt4-dev libqtcore4 libqtgui4 libqtwebkit-dev g++
       cd liteide/build
       ./update_pkg.sh
       export QTDIR=/usr
       ./build_linux.sh
       cd ~/liteide/liteidex
       ./linux_deploy.sh
       cd ~/liteide/liteidex/liteide/bin 
       ./liteide
   ```

5. Open ipsd with liteide 

6. Select the platform you needed, modify current environment according to step 1 and 3
    Modify GOROOT and PATH

7. Compile->Build

## Build Working Environment

IPSD must work with IPSC , you need to do as following to create a environment to create static site.

1. Download IPSD

2. Download IPSC

3. Unzip IPSD

4. Unzip IPSC

5. Copy all the files to IPSD folder

## How it works

There are a original source folder which contains source files, it looks like

* Files    // The files contained in this folder will be added to IPSC as normal file

  * SubFolder1
    * SubFile1
    * SubFile2
  * SubFolder2
  * File1
  * File2

* Html  //The files contained in this folder will be added to ipsc as HTML file

  * H1.html
  * H1.png  //xx.png will be used as title iamge for xx.html
  * H2.html
  * H2.png

* Link 

  * Link.txt //Links in this file will be added to IPSC as hyperlink
  * L1.png //xx.png will be used as title image for link xx
  * L2.png

* Markdown //Markdown files in this folder will be added to IPSC as markdown file

  * M1.md 
  * M1.png //xx.png will be used as title image for xx.md
  * M2.md
  * M2.png

* Templates

  * Blank.md
  * News.md

* monitor.sm  //definition file

  

You can add/remove/update  md,html and link in this folder,  add/remove/update their meta data or title iamge. Before that, you need to connect this folder with a ipsc site project. After you modified the orginial folder,run call ipsd to run monitor, ipsd will check this folder and send all the changes(Add/Remove/Update) to IPSC adn call ipsc to compile the site again. 



## Commands

* New Monitor 

```bash
ipsd -Command NewMonitor -SiteFolder -SiteTitle  -MonitorFolder
```

Create a new monitor, connect a originail source folder to a ipsc site project folder

Example:

```bash
ipsd -Command NewMonitor -SiteFolder "F:\TestSite" -SiteTitle "Test Site" -MonitorFolder "F:\WatchdogSpace"
```

* Run Monitor

```bash
ipsd -Command RunMonitor -MonitorFolder -IndexPageSize
```

Run the monitor defined in MonitorFolder , if there are any change in the monitor folder (add delete or update), will update the changes to ispc and then compile site with IndexPageSize

IndexPageSize (for index page and more page of site, for more information, read QuickHelp.txt of FullHelp.txt of ipsc)

*  Normal 	index(more) page will contain 20 items

* Small  	index(more) page will contain 10 items

*  VerySmall	index(more) page will contain 5  items

*  Big		index(more) page will contain 30 items

  

Example:

```bash
ipsd -Command RunMonitor -MonitorFolder "F:\WatchdogSpace" -IndexPageSize "VerySmall"
```

* List Normal File
	
```bash
ipsd -Command ListNormalFile -MonitorFolder 
```

List all the normal files that already added to the connected site project
	
Example:

```bash
ipsd -Command ListNormalFile -MonitorFolder "F:\WatchdogSpace"
```

## Build a Tool-Chain to create and publish site to IPFS

1. Build IPSD working environment as "Build Working Environment"
2. Download IPSP and unzip it
3. Use IPSC to create a ipsc site project, create ipsc source folder and ipsc output folder
  ```bash
    ipsc -Command "NewSite" -SiteFolder "F:\TestSite" -SiteTitle "Test Site" -SiteAuthor "Chao(sdxianchao@gmail.com)" -SiteDescription "Test Site for IPSC" -OutputFolder "F:\SiteOutputFolder"
  ```
4. Create a empty original source folder for IPSD. Now there are three folders
  * Orignial Folder F:\WatchdogSpace which contains the orginal files for the site
  * Site Source Folder F:\TestSite which contains files for the ipsc site
  * Site Output Folder F:\SiteOutputFolder which contains generated files of site
5. Connect Original source folder with IPSC by ipsd
```bash
  ipsd -Command NewMonitor -SiteFolder "F:\TestSite" -SiteTitle "Test Site" -MonitorFolder "F:\WatchdogSpace"
```
6. Use IPSP to monitor IPSC output folder
```bash
  ipsp -SiteFolder "F:\TestSite" -MonitorInterval 600
```

## Build a Tool-Chain to create and publish site to normal web server (IIS/Apache)

1. Build IPSD working environment as "Build Working Environment"
2. Use IPSC to create a ipsc site project, create ipsc source folder and ipsc output folder
  ```bash
    ipsc -Command "NewSite" -SiteFolder "F:\TestSite" -SiteTitle "Test Site" -SiteAuthor "Chao(sdxianchao@gmail.com)" -SiteDescription "Test Site for IPSC" -OutputFolder "F:\SiteOutputFolder"
  ```
3. Add ipsc output folder to web server
4. Create a empty original source folder for IPSD
 * Orignial Folder F:\WatchdogSpace which contains the orginal files for the site
  * Site Source Folder F:\TestSite which contains files for the ipsc site
  * Site Output Folder F:\SiteOutputFolder which contains generated files of site
5. Connect Original source folder with IPSC by ipsd
```bash
  ipsd -Command NewMonitor -SiteFolder "F:\TestSite" -SiteTitle "Test Site" -MonitorFolder "F:\WatchdogSpace"
```

## Use the Tool-**Chain**

#### 1. Add File 

* #### Markdown File

1. Open original source folder->Templates Folder, copy a md file as template

2. Open it with Markdown Editor

3. Edit the file

4. Modify Meta data as follows (DON'T change METADATA_END_29b43fcf-5b71-4b15-a048-46765f5ef048)

   [//]: #  "Title : Download ipsc"
   [//]: #  "Author: xxxx"
   [//]: #  "Description: Download ipsc "
   [//]: #  "IsTop:True"
   [//]: # "METADATA_END_29b43fcf-5b71-4b15-a048-46765f5ef048"

5. Save the file

6. Create a png image with the same name, and its size should be smaller than 30 KB

7. Copy md and png file to the Original Source Folder->Markdown



* #### Html File

1.  Create a new Html File

2. Modify the meta data 

   <!--
   [//]: #  "Title : Download ipsd"
   [//]: #  "Author: xxx"
   [//]: #  "Description:Download ipsd"
   [//]: #  "IsTop:true"
   -->

3. Save the file

4. Create a png image with the same name, and its size should be smaller than 30 KB

5. Copy md and png file to the Original Source Folder->Html



* #### Link

  1. 

* #### Normal File

  1. Copy a file or folder to Originial Source Folder->Files



### 2. Run Monitor

Use ipsd runmonitor  to update the changes in original source folder to ipsc then call ipsc to compile the site again. ipsc output folder will changed, and as the ipsp is monitoring the output folder of ipsc,  it will detect the change of ipsc output folder and publish the site in ipsc output folder to ipfs agian.

## Raise A Issue

Send email to sdxianchao@gmail.com 



## Maintainers

[@aStarProgrammer](https://github.com/aStarProgrammer).


## License

[MIT](LICENSE)