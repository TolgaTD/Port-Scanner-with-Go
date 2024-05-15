# Port-Scanner-with-Go
Fast Port Scanner
This project is a Go application designed to quickly scan for open ports on a specified target address. Users can choose to scan a range of ports, specific ports, or predefined sets of ports. Additionally, scan results can optionally be saved to a file.

Features
Various port scanning options:
Scan a specified range of ports
Scan specific ports
Scan the first 100 ports
Scan the first 1000 ports
Scan all ports
Option to save results to a file
User interface 

Installation
To use this project, Go must be installed on your computer. You can download and install the latest version of Go from Go's official website.

Clone the project to your local machine using the following Git command:

•	git clone https://github.com/yourusername/projectname.git

•	cd projectname

Usage
To run the program, enter the following command in your terminal:

•	go run main.go

After running the program, follow the on-screen prompts to specify the target address and scanning options.





Example Usage
If you enter the target address 192.168.1.1 and the port range 1-100, the program will scan ports 1 through 100 on this IP address.

Configuration
To adjust the scanning speed and other parameters, you can modify the rate limiter settings in the main.go file.

License
This project is licensed under the MIT License.

