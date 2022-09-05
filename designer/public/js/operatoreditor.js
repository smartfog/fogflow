function registerAllBlocks(blocks){

blocks.register({
    name: "Operator",
    description: "To specify an operator",
    fields: [
        {
            name: "Name",
            type: "string",
            defaultValue: "unknown",
            attrs: "editable"
        },     
        {
            name: "Description",
            type: "longtext",
            defaultValue: "fogflow",
            attrs: "editable"
        },    
        {
            name: "Parameters",
            type: "Parameter",
            attrs: "input"
        }
    ]
});

blocks.register({
    name: "Parameter",
    description: "To specify an controllable parameter of the operator",
    fields: [
        {
            name: "Name",
            type: "string",
            choices: ["service_port", "virtual_network_interface"],            
            attrs: "editable"
        },{
            name: "Value",
            type: "string",
            attrs: "editable"
        },
        {
            name: "Parameter",
            type: "Parameter",
            attrs: "output"
        }        
    ]
});


}



