Steps to execute the Unittest script and generate the test report

1. Change the append path of  in the script test_allsubscription.py
	for example : script already have path sys.path.append('/home/ubuntu/fogflow/test/UnitTest')

2. Change the broker port and ip according to running broker port and ip.
	for example:  script already have brokerIp="http://192.168.100.120:8070"

3. run the command : py.test 

4. To generate the report run the following command
	pytest --html=report.html
	The above command will generate  two file one report.html and second assets/style.css
	for more information follow link : https://github.com/pytest-dev/pytest-html

   
