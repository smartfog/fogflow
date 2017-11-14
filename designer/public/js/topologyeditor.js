
function registerAllBlocks(blocks, operators){

console.log("operator list: ", operators);

blocks.register({
    name: "Task",
    description: "To specify a data processing task",
    fields: [
        {
            name: "name",
            type: "string",
            attrs: "editable"
        },    
        {
            name: "operator",
            type: "choice",            
            choices: operators,
            attrs: "editable"
        },    
        {
            name: "groupby",
            type: "string",
            defaultValue: "all",
            attrs: "editable"
        }, 
        {
            name: "inputs",
            label: "shuffling of input streams",
            type: "choice[]",
            hide: true,
            choices: ["unicast","broadcast"],
            defaultValue: ["unicast"],       
            attrs: "editable input",
            card: '1',                            
            dynamicLabel: function(block, x) {
                return block.getValue('inputs')[x];
            }             
        },                        
        {
            name: "outputs",
            label: "entity type of output streams",
            type: "string[]",
            defaultValue: ["Out"],                                   
            hide: true,
            attrs: "editable output",
            dynamicLabel: function(block, x) {
                return block.getValue('outputs')[x];
            }                                  
        }      
    ]
});

blocks.register({
    name: "InputStream",
    description: "To define the input stream",
    fields: [
        {
            name: "EntityType",
            type: "string",
            attrs: "editable"
        },
        {
            name: "scoped",
            type: "bool",
            defaultValue: false,
            attrs: "editable"
        },
        {
            name: "stream",
            label: "Stream",
            attrs: "output",
            type: "string"
        }                          
    ]
});


}



