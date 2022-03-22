*****************************************
Fogflow Automation Testing
*****************************************
Automation testing is a software testing technique that executes a test case suite using particular automated testing software tools. Please follow the following steps to deploy the entire FogFlow system on a single Linux machine before your test. 

`Set up all FogFlow component on a single machine`_

.. _`Set up all FogFlow component on a single machine`: https://fogflow.readthedocs.io/en/latest/setup.html

Once the FogFlow is up and running, an end-to-end function test can be carried out by using python library pytest with a test plan. 
More detailed steps are available below. 

`Getting started with Pytest`_

.. _`Getting started with Pytest`: https://docs.pytest.org/en/6.2.x/getting-started.html?msclkid=a752f7b6a8d711ec92b8b7afabf6eda5

NGSI-V1 Unit Testing
===========================================================

`NGSI-V1 Test Suite`_

.. _`NGSI-V1 Test Suite`: https://github.com/smartfog/fogflow/tree/development/test/UnitTest/v1

You should get an output similar to the following one:

.. code-block:: console
    
    root@fog-3:~/usecase_fogflow/test/UnitTest/v1# pytest -v test_casesNGSIv1.py
    ============================================================= test session starts ==============================================================
    platform linux2 -- Python 2.7.17, pytest-3.3.2, py-1.5.2, pluggy-0.6.0 -- /usr/bin/python2
    cachedir: .cache
    rootdir: /root/usecase_fogflow/test/UnitTest/v1, inifile:
    collected 37 items

    test_casesNGSIv1.py::test_getSubscription1 PASSED                     [  2%]
    test_casesNGSIv1.py::test_getSubscription2 PASSED                     [  5%]
    test_casesNGSIv1.py::test_getSubscription3 PASSED                     [  8%]
    test_casesNGSIv1.py::test_getSubscription4 PASSED                     [ 10%]
    test_casesNGSIv1.py::test_getSubscription5 PASSED                     [ 13%]
    test_casesNGSIv1.py::test_getSubscription6 PASSED                     [ 16%]
    test_casesNGSIv1.py::test_getSubscription7 PASSED                     [ 18%]
    test_casesNGSIv1.py::test_getSubscription8 PASSED                     [ 21%]
    test_casesNGSIv1.py::test_getSubscription9 PASSED                     [ 24%]
    test_casesNGSIv1.py::test_getSubscription10 PASSED                    [ 27%]

        .. note:: no. of test cases are more than the no. of test cases shown above for example. 


NGSI-LD Unit Testing 
===========================================================

`NGSI-LD Test Suite`_

.. _`NGSI-LD Test Suite`: https://github.com/smartfog/fogflow/tree/development/test/UnitTest/NGSI-LD

You should get an output similar to the following one:

.. code-block:: console
    
    root@fog-3:~/usecase_fogflow/test/UnitTest/NGSI-LD# pytest -v test_casesNGSI-LD.py
    ============================================================= test session starts ==============================================================
    platform linux2 -- Python 2.7.17, pytest-3.3.2, py-1.5.2, pluggy-0.6.0 -- /usr/bin/python2
    cachedir: .cache
    rootdir: /root/usecase_fogflow/test/UnitTest/NGSI-LD, inifile:
    collected 156 items

    test_casesNGSI-LD.py::test_case1 PASSED                       [  0%]
    test_casesNGSI-LD.py::test_case2 PASSED                       [  1%]
    test_casesNGSI-LD.py::test_case4 PASSED                       [  1%]
    test_casesNGSI-LD.py::test_case5 PASSED                       [  2%]
    test_casesNGSI-LD.py::test_case6 PASSED                       [  3%]
    test_casesNGSI-LD.py::test_case7 PASSED                       [  3%]
    test_casesNGSI-LD.py::test_case8 PASSED                       [  4%]
    test_casesNGSI-LD.py::test_case9 PASSED                       [  5%]
    test_casesNGSI-LD.py::test_case10 PASSED                      [  5%]
    test_casesNGSI-LD.py::test_case11 PASSED                      [  6%]
    test_casesNGSI-LD.py::test_case12 PASSED                      [  7%]
    test_casesNGSI-LD.py::test_case13 PASSED                      [  7%]
    test_casesNGSI-LD.py::test_case14 PASSED                      [  8%]
    test_casesNGSI-LD.py::test_case15 PASSED                      [  8%]
    test_casesNGSI-LD.py::test_case16 PASSED                      [  9%]
    test_casesNGSI-LD.py::test_case17 PASSED                      [ 10%]
    test_casesNGSI-LD.py::test_case18 PASSED                      [ 10%]
    test_casesNGSI-LD.py::test_case19 PASSED                      [ 11%]
    test_casesNGSI-LD.py::test_case20 PASSED                      [ 12%]

        .. note:: no. of test cases are more than the no. of test cases shown above for example. 


Persistance Unit Testing 
===========================================================

`Persistance Test Suite`_

.. _`Persistance Test Suite`: https://github.com/smartfog/fogflow/tree/development/test/UnitTest/persistance

You should get an output similar to the following one:

.. code-block:: console

    root@fog-3:~/usecase_fogflow/test/UnitTest/persistance# pytest -v test_persistance.py
    ============================================================= test session starts ==============================================================
    platform linux2 -- Python 2.7.17, pytest-3.3.2, py-1.5.2, pluggy-0.6.0 -- /usr/bin/python2
    cachedir: .cache
    rootdir: /root/usecase_fogflow/test/UnitTest/persistance, inifile:
    collected 21 items

    test_persistance.py::test_persistOPerator PASSED                                 [  4%]
    test_persistance.py::test_persistOPerator1 PASSED                                [  9%]
    test_persistance.py::test_persistOPerator2 PASSED                                [ 14%]
    test_persistance.py::test_persistOPerator3 PASSED                                [ 19%]
    test_persistance.py::test_persistFogFunction PASSED                              [ 23%]
    test_persistance.py::test_persistFogFunction1 PASSED                             [ 28%]
    test_persistance.py::test_persistFogFunction2 PASSED                             [ 33%]
    test_persistance.py::test_persistFogFunction3 PASSED                             [ 38%]
    test_persistance.py::test_persistFogFunction4 PASSED                             [ 42%]
    test_persistance.py::test_persistDockerImage PASSED                              [ 47%]
    test_persistance.py::test_persistDockerImage1 PASSED                             [ 52%]
    test_persistance.py::test_persistDockerImage2 PASSED                             [ 57%]
    test_persistance.py::test_persistDockerImage3 PASSED                             [ 61%]
    test_persistance.py::test_persistDockerImage4 PASSED                             [ 66%]
    test_persistance.py::test_persistopology PASSED                                  [ 71%]
    test_persistance.py::test_persistopology1 PASSED                                 [ 76%]
    test_persistance.py::test_persistopology2 PASSED                                 [ 80%]
    test_persistance.py::test_persistintent PASSED                                   [ 85%]
    test_persistance.py::test_persistintent1 PASSED                                  [ 90%]
    test_persistance.py::test_persistintent2 PASSED                                  [ 95%]
    test_persistance.py::test_persistintent3 PASSED                                  [100%]
