$(function() {

    // initialization  
    var handlers = {};

    addMenuItem('Operator', 'Operator', showOperator);
    addMenuItem('DockerImage', 'Docker Image', showDockerImage);
    
    initOperatorList();
    initDockerImageList();    

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
        
        fetch('/operator').then(res => res.json()).then(operators => {
            var opList = Object.values(operators)
            displayOperatorList(opList)
        });        
    }

    function initOperatorList(){
        fetch('/operator').then(res => res.json()).then(opList => {
            if (Object.keys(opList).length === 0) {
                var operators = defaultOperatorList();
                fetch("/operator", {
                    method: "POST",
                    headers: {
                        Accept: "application/json",
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(operators)
                })
                .then(response => {
                    console.log("send the initial list of operators", response.status)
                    showOperator();
                })
                .catch(err => console.log(err));
            } else {
                showOperator();
            }               
        })
    }

    function queryOperatorList() {
        fetch('/operator').then(res => res.json()).then(operators => {
            Object.values(operators).forEach(operator => {
                var option = document.createElement("option");
                option.text = operator.name; 
                
                var operatorList = document.getElementById("OperatorList");
                operatorList.add(option);                
            });           
        }); 
    }

    function displayOperatorList(operators) {
        if (operators == undefined || operators.length == 0) {
            $('#operatorList').html('');
            return
        }

        var html = '<table class="table table-striped table-bordered table-condensed">';

        html += '<thead><tr>';
        html += '<th>Operator</th>';
        html += '<th>Description</th>';
        html += '<th>#Parameters</th>';
        html += '</tr></thead>';

        for (var i = 0; i < operators.length; i++) {
            html += '<tr>';
            html += '<td>' + operators[i].name + '</td>';
            html += '<td>' + operators[i].description + '</td>';
            
            if ('parameters' in operators[i])
                html += '<td>' + operators[i].parameters.length + '</td>';
            else
                html += '<td>' + 0 + '</td>';

            html += '</tr>';
        }

        html += '</table>';

        $('#operatorList').html(html);
    }
    
    
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
        var operator = boardScene2Operator(scene);
        
        if(operator.name && operator.name != "unknown")
            submitOperator(operator, scene);
        else
            alert('please provide the required inputs');
    }

    function submitOperator(operator) {
        console.log(operator);
        
        fetch("/operator", {
            method: "POST",
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json"
            },
            body: JSON.stringify([operator])
        })
        .then(response => {
            console.log("send the initial list of operators: ", response.status)
            showOperator();
        })
        .catch(err => console.log(err));
    }

    function boardScene2Operator(scene) {
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
                        parameter.value = block.values.value;

                        parameters.push(parameter);
                    }
                }
            }
        }

        return parameters;
    }

    function showDockerImage() {
        $('#info').html('list of docker images in the docker registry');

        var html = '<div style="margin-bottom: 10px;"><button id="registerDockerImage" type="button" class="btn btn-primary">register</button></div>';
        html += '<div id="dockerImageList"></div>';

        $('#content').html(html);

        updateDockerImageList();

        $("#registerDockerImage").click(function() {
            dockerImageRegistration();
        });
    }

    function initDockerImageList(){                        
        fetch('/dockerimage').then(res => res.json()).then(imageList => {
            if (Object.keys(imageList).length === 0) {
                var images = defaultDockerImageList();
                fetch("/dockerimage", {
                    method: "POST",
                    headers: {
                        Accept: "application/json",
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(images)
                })
                .then(response => {
                    console.log("send the initial list of docker images: ", response.status)
                })
                .catch(err => console.log(err));
            }               
        })          
    }
    
//    function addDockerImage(image) {
        
//        //register a new docker image        
//        fetch("/dockerimage", {
//            method: "POST",
//            headers: {
//                Accept: "application/json",
//                "Content-Type": "application/json"
//            },
//            body: JSON.stringify([image])
//        })
//        .then(response => {
//            console.log("add a new docker image: ", response.status)
//        })
//        .catch(err => console.log(err));        
//    }

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

        var tag = $('#imageTag').val();
        if (tag == '') {
            tag = 'latest';
        }

        var hwType = $('#hwType option:selected').val();
        var osType = $('#osType option:selected').val();
        var operatorName = $('#OperatorList option:selected').val();

        var prefetched = document.getElementById('Prefetched').checked;

        if (image == '' || tag == '' || hwType == '' || osType == '' || operatorName == '') {
            alert('please provide the required inputs');
            return;
        }

        //register a new docker image
        var newImageObject = {};
        newImageObject.name = image
        newImageObject.hwType = hwType
        newImageObject.osType = osType
        newImageObject.operatorName = operatorName
        newImageObject.prefetched = prefetched
        newImageObject.tag = tag
        
        console.log([newImageObject]);        
        
        fetch("/dockerimage", {
            method: "POST",
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json"
            },
            body: JSON.stringify([newImageObject])
        })
        .then(data => {
            showDockerImage();            
        })
        .catch(err => {
            console.log('failed to register the new docker image, ', err);
        });      
    }


    function updateDockerImageList() {
        fetch('/dockerimage').then(res => res.json()).then(dockerimages => {
            var imageList = Object.values(dockerimages);
            displayDockerImageList(imageList);
        })
        .catch(err => {
            console.log('failed to fetch the list of docker images, ', err);
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
