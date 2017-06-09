# Image fetcher from danbooru
This small Go script will fetch a random picture from Danbooru within the set interval and save it in a directory while deleting the old one (so there will be at most only 1 image in the directory at a time). You can set your wallpaper setting to "Slideshow", point to that folder and synchronize the duration so that each new picture is set to be next wallpaper.  
There are 3 parameters that can be changed:
- Danbooru API endpoint (and effectively, the tags you want)
- The local folder to save images
- The interval between each download

You can further change the criteria for a qualified image (minimum height/width, width-to-height ratio) in the script.  
Usage: edit the constants (and whatever you want), then run `go run fetcher.go`. Or you can compile it beforehand with `go build fetcher.go` to generate an executable.