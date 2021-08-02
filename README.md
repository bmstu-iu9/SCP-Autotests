# SCP-Autotests
Разработка автотестов для суперкомпиляторов MSCP-A и SCP4.

---

### This application can
*    Convert a refal program with parametrized entry point into compiled
*    Run the supercompilers on these programs and catch errors
*    Compare default and residual programs outputs
*    Compare times of work of different supercompilers on each test

---

### For launch autotests
_You need the refal-5's or refal-5-lambda's compiler to be installed_
1. add the MSCPAver _{your version's number}_ to the project's directory
2. **_go build_**
3. **_./autotests_**
    * **-scp** {_version's number_}
    * **-v** {_refal version_}
    * **-path** {_path to json tests storage_}
    * **-rsd** {_to save rsd files after program's ending_}
    * **-time** {_to make time of this scp standard_}

#### default parameters:
* scp - _1_
* v - _default_
* rsd - _false_
* time - _false_

_You can add your tests (refal programs) into tests/ and your supercompilers into main directory_
