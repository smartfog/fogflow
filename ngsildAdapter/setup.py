import sys, os
testdir = os.path.dirname(__file__)
srcdir = '../module'
sys.path.insert(0, os.path.abspath(os.path.join(testdir, srcdir)))
#from data_model.ld_generate  import ngsi_data_creation
sys.path.append('/root/TRANSFORMER/Next_transform/fogflow/ngsildAdapter/module')
