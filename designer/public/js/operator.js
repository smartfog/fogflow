$(function() {

    // initialization  
    var handlers = {};

    //connect to the broker
    var client = new NGSI10Client(config.brokerURL);

    // to interact with designer for internal fogflow entities
    var clientDes = new NGSI10Client('./internal');

    console.log(config.brokerURL);

    addMenuItem('Operator', 'Operator', showOperator);
    addMenuItem('DockerImage', 'Docker Image', showDockerImage);
    initOperatorList();
    initDockerImageList();
    showOperator();


    $(window).on('hashchange', function() {
        var hash = window.location.hash;
        selectMenuItem(location.hash.substring(1));
    });

    function addMenuItem(id, name, func) {
        handlers[id] = func;
        $('#menu').append('<li id="' + id + '"><a href="' + '#' + id + '">' + name + '</a></li>');
    }

    function selectMenuItem(name) {
        $('#menu li').removeClass('active');
        var element = $('#' + name);
        element.addClass('active');

        var handler = handlers[name];
        if (handler != null) {
            handler();
        }
    }


    function showOperator() {
        $('#info').html('list of all registered operators');

        var html = '<div style="margin-bottom: 10px;"><button id="registerOperator" type="button" class="btn btn-primary">register</button></div>';
        html += '<div id="operatorList"></div>';

        $('#content').html(html);

        $("#registerOperator").click(function() {
            showOperatorEditor();
        });

        var queryReq = {}
        queryReq = { internalType: "Operator", updateAction: "UPDATE" };
        clientDes.getContext(queryReq).then(function(operators) {
            console.log("show operator ",operators);
            displayOperatorList(operators)
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });
    }

    function initOperatorList(){
        var queryReq = {}
        queryReq = { internalType: "Operator", updateAction: "UPDATE" };
        clientDes.getContext(queryReq).then(function(opList) {
            console.log("inside init operator list ",opList);
            if (opList.data.length == 0) {
                console.log("inside init operator list ******* ",opList);
                for (var i = 0; i < defaultOperatorList().length; i++) { 
                    opObj = defaultOperatorList()[i];
                    var operator ={} 
                    var attribute = opObj;
                    operator.attribute = attribute;
                    operator.updateAction = "UPDATE"
                    operator.internalType = "Operator"
                    submitOperator(operator,"")
                }
            }
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query fog functions');
        });
    }

    function queryOperatorList() {
        var queryReq = {}
        queryReq = { internalType: "Operator", updateAction: "UPDATE" };
        clientDes.getContext(queryReq).then(function(operators) {
            for (var i = 0; i < operators.data.length; i++) {
                if (Object.keys(operators.data[i]).length === 0) continue;
                var option = document.createElement("option");
                option.text = operators.data[i].name; 
                var operatorList = document.getElementById("OperatorList");
                operatorList.add(option);
            }

            // add it into the select list        
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });
    }

    //vinod: start
    function displayOperatorList(operators) {
        operators = operators.data
        console.log("new operator list",operators)
        if (operators == undefined || operators.length == 0) {
            $('#operatorList').html('');
            return
        }

        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>Operator</th>';
        html += '<th>Name</th>';
        html += '<th>Description</th>';
        html += '<th>#Parameters</th>';
        html += '<th>#Implementations</th>';
        html += '</tr></thead>';

        for (var i = 0; i < operators.length; i++) {
            if (Object.keys(operators[i]).length === 0) continue;
            console.log("for loop ",operators[i].name)

            html += '<tr>';
            html += '<td>' + operators[i].name + '</td>';
            html += '<td>' + operators[i].name + '</td>';
            html += '<td>' + operators[i].description + '</td>';
            html += '<td>' + operators[i].parameters.length + '</td>';
            html += '<td>' + 0 + '</td>';

            html += '</tr>';
        }

        html += '</table>';

        $('#operatorList').html(html);
    }// vinod:end

    function showOperatorEditor() {
        $('#info').html('to specify an operator');

        var html = '';

        html += '<div><button id="generateOperator" type="button" class="btn btn-primary">Submit</button></div>';
        html += '<div id="blocks" style="width:1000px; margin-top: 5px; height:400px"></div>';

        $('#content').html(html);

        blocks = new Blocks();

        registerAllBlocks(blocks);

        blocks.run('#blocks');

        blocks.types.addCompatibility('string', 'choice');

        blocks.ready(function() {
            // associate functions to clickable buttons
            $('#generateOperator').click(function() {
                generateOperator(blocks.export());
            });
        });
    }


    function generateOperator(scene) {
        // construct the operator based on the design board
        var operator ={} 
        var attribute = boardScene2Operator(scene);
        operator.attribute = attribute;
        operator.updateAction = "UPDATE"
        operator.internalType = "Operator"
        console.log("------------------op gen---",operator);

        // submit this operator
        
        if(operator.attribute && operator.attribute.name && operator.attribute.name != "unknown")
        submitOperator(operator, scene);
        else {
            alert('please provide the required inputs');
            return;
        }
    }

    function submitOperator(operator, designboard) {
        clientDes.updateContext(operator).then(function(data) {
            showOperator();
        }).catch(function(error) {
            console.log('failed to submit the defined operator');
        });
    }

    function boardScene2Operator(scene) {
        console.log(scene);
        var operator = {};

        for (var i = 0; i < scene.blocks.length; i++) {
            var block = scene.blocks[i];

            if (block.type == "Operator") {
                operator.name = block.values['name'];
                operator.description = block.values['description'];

                // construct its controllable parameters
                operator.parameters = findInputParameters(scene, block.id);

                break;
            }
        }

        console.log(operator);
        return operator;
    }


    function findInputParameters(scene, blockId) {
        var parameters = [];

        for (var i = 0; i < scene.edges.length; i++) {
            var edge = scene.edges[i];

            if (edge.block2 == blockId) {
                for (var j = 0; j < scene.blocks.length; j++) {
                    var block = scene.blocks[j];

                    if (block.id == edge.block1) {
                        var parameter = {};
                        parameter.name = block.values.name
                        parameter.values = block.values.values;

                        parameters.push(parameter);
                    }
                }
            }
        }

        return parameters;
    }

    function showDockerImage() {
        console.log("show docker +++++")
        $('#info').html('list of docker images in the docker registry');

        var html = '<div style="margin-bottom: 10px;"><button id="registerDockerImage" type="button" class="btn btn-primary">register</button></div>';
        html += '<div id="dockerImageList"></div>';

        $('#content').html(html);

        updateDockerImageList();

        $("#registerDockerImage").click(function() {
            dockerImageRegistration();
        });
    }

    
    
    function defaultOperatorList(){
        var operatorList = [{
            name: "nodejs",
            description: "",
            parameters: []
        }, {
            name: "python",
            description: "",
            parameters: []
        }, {
            name: "iotagent",
            description: "",
            parameters: []
        }, {
            name: "counter",
            description: "",
            parameters: []
        }, {
            name: "anomaly",
            description: "",
            parameters: []
        }, {
            name: "facefinder",
            description: "",
            parameters: []
        }, {
            name: "connectedcar",
            description: "",
            parameters: []
        }, {
            name: "recommender",
            description: "",
            parameters: []
        }, {
            name: "privatesite",
            description: "",
            parameters: []
        }, {
            name: "publicsite",
            description: "",
            parameters: []
        }, {
            name: "pushbutton",
            description: "",
            parameters: []
        }, {
            name: "acoustic",
            description: "",
            parameters: []
        }, {
            name: "speaker",
            description: "",
            parameters: []
        }, {
            name: "dummy",
            description: "",
            parameters: []
        }, {
            name: "geohash",
            description: "",
            parameters: []
        }, {
            name: "converter",
            description: "",
            parameters: []
        }, {
            name: "predictor",
            description: "",
            parameters: []
        }, {
            name: "controller",
            description: "",
            parameters: []
        }, {
            name: "detector",
            description: "",
            parameters: []
        }, {
            name: "LDanomaly",
            description: "",
            parameters: []
        }, {
            name: "LDCounter",
            description: "",
            parameters: []
        }];
       
        return operatorList
    }

  

    function defaultDockerImageList() {
        var imageList = [{
            name: "fogflow/nodejs",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "nodejs",
            prefetched: true
        }, {
            name: "fogflow/python",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "python",
            prefetched: false
        }, {
            name: "fogflow/counter",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "counter",
            prefetched: false
        }, {
            name: "fogflow/anomaly",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "anomaly",
            prefetched: false
        }, {
            name: "fogflow/facefinder",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "facefinder",
            prefetched: false
        }, {
            name: "fogflow/connectedcar",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "connectedcar",
            prefetched: false
        }, {
            name: "fiware/iotagent-json",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "iotagent",
            prefetched: false
        }, {
            name: "fogflow/recommender",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "recommender",
            prefetched: false
        }, {
            name: "fogflow/privatesite",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "privatesite",
            prefetched: false
        }, {
            name: "fogflow/publicsite",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "publicsite",
            prefetched: false
        }, {
            name: "pushbutton",
            tag: "latest",
            hwType: "ARM",
            osType: "Linux",
            operatorName: "pushbutton",
            prefetched: false
        }, {
            name: "acoustic",
            tag: "latest",
            hwType: "ARM",
            osType: "Linux",
            operatorName: "acoustic",
            prefetched: false
        }, {
            name: "speaker",
            tag: "latest",
            hwType: "ARM",
            osType: "Linux",
            operatorName: "speaker",
            prefetched: false
        }, {
            name: "pushbutton",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "pushbutton",
            prefetched: false
        }, {
            name: "acoustic",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "acoustic",
            prefetched: false
        }, {
            name: "speaker",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "speaker",
            prefetched: false
        }, {
            name: "fogflow/dummy",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "dummy",
            prefetched: false
        }, {
            name: "geohash",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "geohash",
            prefetched: false
        }, {
            name: "converter",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "converter",
            prefetched: false
        }, {
            name: "predictor",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "predictor",
            prefetched: false
        }, {
            name: "controller",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "controller",
            prefetched: false
        }, {
            name: "detector",
            tag: "latest",
            hwType: "ARM",
            osType: "Linux",
            operatorName: "detector",
            prefetched: false
        }, {
            name: "fogflow/ldanomaly",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "LDanomaly",
            prefetched: false
        }, {
            name: "fogflow/ldcounter",
            tag: "latest",
            hwType: "X86",
            osType: "Linux",
            operatorName: "LDCounter",
            prefetched: false
        }];

        return imageList;
    }

    function initDockerImageList(){
        var queryReq = {}
        queryReq = { internalType: "DockerImage", updateAction: "UPDATE" };
        clientDes.getContext(queryReq).then(function(doList) {
            if (doList.data.length == 0) {
                for (var i = 0; i < defaultDockerImageList().length; i++) { 
                    dockerObj = defaultDockerImageList()[i];
                    var newImageObject = {};
                    var attribute = {};
                    attribute.name = dockerObj.name;
                    attribute.hwType = dockerObj.hwType;
                    attribute.osType = dockerObj.osType;
                    attribute.operatorName = dockerObj.operatorName;
                    attribute.prefetched = dockerObj.prefetched;
                    attribute.tag = dockerObj.tag;
                    newImageObject.attribute = attribute;
                    newImageObject.internalType = "DockerImage"
                    newImageObject.updateAction = "UPDATE"

                    clientDes.updateContext(newImageObject).then(function(data) {
                        console.log(data);
                    }).catch(function(error) {
                        console.log('failed to register the new device object');
                    });
                }
            }
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query fog functions');
        });
    }
    function addDockerImage(image) {
        //register a new docker image
        var newImageObject = {};

        newImageObject.entityId = {
            id: image.name + '.' + image.tag,
            type: 'DockerImage',
            isPattern: false
        };

        newImageObject.attributes = {};
        newImageObject.attributes.image = { type: 'string', value: image.name };
        newImageObject.attributes.tag = { type: 'string', value: image.tag };
        newImageObject.attributes.hwType = { type: 'string', value: image.hwType };
        newImageObject.attributes.osType = { type: 'string', value: image.osType };
        newImageObject.attributes.operator = { type: 'string', value: image.operatorName };
        newImageObject.attributes.prefetched = { type: 'boolean', value: image.prefetched };

        newImageObject.metadata = {};
        newImageObject.metadata.operator = {
            type: 'string',
            value: image.operatorName
        };

        var geoScope = {};
        geoScope.type = "global"
        geoScope.value = "global"
        newImageObject.metadata.location = geoScope;

        clientDes.updateContext(newImageObject).then(function(data) {
            console.log(data);
        }).catch(function(error) {
            console.log('failed to register the new device object');
        });

    }

    function dockerImageRegistration() {
        $('#info').html('New docker image registration');

        var html = '<div id="dockerRegistration" class="form-horizontal"><fieldset>';

        html += '<div class="control-group"><label class="control-label" for="input01">Image(*)</label>';
        html += '<div class="controls"><input type="text" class="input-xlarge" id="dockerImageName">';
        html += '</div></div>';

        html += '<div class="control-group"><label class="control-label" for="input01">Tag(*)</label>';
        html += '<div class="controls"><input type="text" class="input-xlarge" id="imageTag" placeholder="latest">';
        html += '</div></div>';

        html += '<div class="control-group"><label class="control-label" for="input01">HardwareType(*)</label>';
        html += '<div class="controls"><select id="hwType"><option>X86</option><option>ARM</option></select></div>'
        html += '</div>';

        html += '<div class="control-group"><label class="control-label" for="input01">OSType(*)</label>';
        html += '<div class="controls"><select id="osType"><option>Linux</option><option>Windows</option></select></div>'
        html += '</div>';

        html += '<div class="control-group"><label class="control-label" for="input01">Operator(*)</label>';
        html += '<div class="controls"><select id="OperatorList"></select>';
        html += '</div></div>';

        html += '<div class="control-group"><label class="control-label" for="optionsCheckbox">Prefetched</label>';
        html += '<div class="controls"> <label class="checkbox"><input type="checkbox" id="Prefetched" value="option1">';
        html += 'docker image must be fetched by the platform in advance';
        html += '</label></div>';
        html += '</div>';


        html += '<div class="control-group"><label class="control-label" for="input01"></label>';
        html += '<div class="controls"><button id="submitRegistration" type="button" class="btn btn-primary">Register</button>';
        html += '</div></div>';

        html += '</fieldset></div>';

        $('#content').html(html);

        queryOperatorList();

        // associate functions to clickable buttons
        $('#submitRegistration').click(registerDockerImage);
    }


    function registerDockerImage() {
        console.log('register a new docker image');

        // take the inputs    
        var image = $('#dockerImageName').val();
        console.log(image);

        var tag = $('#imageTag').val();
        if (tag == '') {
            tag = 'latest';
        }

        console.log(tag);

        var hwType = $('#hwType option:selected').val();
        console.log(hwType);

        var osType = $('#osType option:selected').val();
        console.log(osType);

        var operatorName = $('#OperatorList option:selected').val();
        console.log(operatorName);

        var prefetched = document.getElementById('Prefetched').checked;
        console.log(prefetched);


        if (image == '' || tag == '' || hwType == '' || osType == '' || operatorName == '') {
            alert('please provide the required inputs');
            return;
        }

        //register a new docker image
        var newImageObject = {};
        var attribute = {};
        attribute.name = image
        attribute.hwType = hwType
        attribute.osType = osType
        attribute.operatorName = operatorName
        attribute.prefetched = prefetched
        attribute.tag = tag
        newImageObject.attribute = attribute
        newImageObject.internalType = "DockerImage"
        newImageObject.updateAction = "UPDATE"
       
        clientDes.updateContext(newImageObject).then(function(data) {
            console.log(data);

            // show the updated image list
            showDockerImage();
        }).catch(function(error) {
            console.log('failed to register the new device object');
        });

    }


    function updateDockerImageList() {
        var queryReq = {}
        queryReq = { internalType: "DockerImage", updateAction: "UPDATE" };

        clientDes.getContext(queryReq).then(function(imageList) {
            imageList = imageList.data
            console.log("get docker image list ",imageList)
            displayDockerImageList(imageList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });
    }

    function displayDockerImageList(images) {
        if (images == null || images.length == 0) {
            $('#dockerImageList').html('');
            return
        }

        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>Operator</th>';
        html += '<th>Image</th>';
        html += '<th>Tag</th>';
        html += '<th>Hardware Type</th>';
        html += '<th>OS Type</th>';
        html += '<th>Prefetched</th>';
        html += '</tr></thead>';

        for (var i = 0; i < images.length; i++) {
            var dockerImage = images[i];
            console.log("test9999999999999999999 ",dockerImage.hasOwnProperty("operatorName")?dockerImage.operatorName:dockerImage.operater);

            html += '<tr>';
            html += '<td>' + dockerImage.operatorName + '</td>';
            html += '<td>' + dockerImage.name + '</td>';
            html += '<td>' + dockerImage.tag + '</td>';
            html += '<td>' + dockerImage.hwType + '</td>';
            html += '<td>' + dockerImage.osType + '</td>';

            if (dockerImage.prefetched == true) {
                html += '<td><font color="red"><b>' + dockerImage.prefetched + '</b></font></td>';
            } else {
                html += '<td>' + dockerImage.prefetched + '</td>';
            }

            html += '</tr>';
        }

        html += '</table>';

        $('#dockerImageList').html(html);
    }

});
