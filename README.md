# eddy.sh

centralized developer onboarding tool. 

## usage 

download proper `eddy` from release page to suit your operating system. for now only windows and mac are supported. 

## purpose 

in every project I joined most of the documentation on how to setup a project is stored in a confluence or as a tribe knowledge. 

The more tools you require, the more complicated it gets to keep everything up to date.

Few people try to actually make it easier and write scripts to setup a workstation, however people are not coming every day, but the codebase and requirements may vary more often. 

That's why I wanted to create `eddy.sh`, so every tool that I need (or someone else if feasible) can be installed automatically, just by providing correct .yml file. 

## future 

In ideal world I want `eddy.sh` to be able to scan your system to find packages and prefered way of installation, so you won't even have to remember what needs to be setup, just run `eddy.sh scan` and receive a .yml you can share with newcomer. 

## present 

If you prepare .yml file, you can setup your newcomer in minutes instead of days, without jumping through many doc files on how to install CMake on windows machine when your frontend developer needs to be able to compile web assembly to make changes for your company. 

## contribute 

if you want to contribute to the project, it is setup in a way where LLMs can easily generate new tool installation. just add a proper file and put codebase as context and you should be good to go. just confirm whether or not a tool you are adding is installed correctly and raise a pr with a tool you need for your project. 

to run the app just run `go run main.go` from the root folder.

if you want to be able to see additional log messages, run `source scripts/enable_debug.sh`