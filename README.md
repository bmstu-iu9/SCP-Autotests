# SCP-Autotests
Разработка автотестов для суперкомпиляторов MSCP-A и SCP4.

This application can
    -Convert a refal program with parametrized entry point into compiled
    -Run the supercompilers on these programs and catch errors
    -Compare default and residual programs outputs
    -Compare times of work of different supercompilers on each test

For launch autotests:
    -add the MSCPAver{your version's number} to the project's directory
    -go build *.go
    -./autotests 
        -scp {version's number} 
        -v {refal version} 
        -path {path to json tests storage} 
        -rsd {yes or no to save rsd files after program's ending} 
        -time {yes or no to make time of this scp standard}

default parameters:
    scp - 1
    v - default
    rsd - no
    time - no

You can add your tests (refal programs) into tests/ and your supercompilers into main directory
