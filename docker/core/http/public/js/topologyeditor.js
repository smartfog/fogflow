
function registerAllBlocks(blocks, operators){

console.log("operator list: ", operators);

blocks.register({
    name: "Task",
    description: "To specify a data processing task",
    fields: [
        {
            name: "Name",
            type: "string",
            attrs: "editable"
        },    
        {
            name: "Operator",
            type: "choice",            
            choices: operators,
            attrs: "editable"
        },    
        {
            name: "Streams",
            type: "Stream",
            attrs: "input"
        },                        
        {
            name: "Outputs",
            type: "string[]",
            defaultValue: ["Out"],                                   
            hide: true,
            attrs: "editable output",
            dynamicLabel: function(block, x) {
                return block.getValue('Outputs')[x];
            }                                  
        }      
    ]
});


blocks.register({
    name: "Shuffle",
    description: "To decide how the selected entity stream should be forwarded to the instances of next task",
    fields: [
        {
            name: "SelectedAttributes",
            type: "string[]",
            defaultValue: ["all"],
            attrs: "editable"
        },          
        {
            name: "Groupby",
            choices: ["ALL", "EntityID", "EntityType", "EntityAttribute"],
            defaultValue: "EntityID",            
            attrs: "editable"
        },              
        {
            name: "In",
            type: "string",
            attrs: "input"
        },        
        {
            name: "Stream",
            type: "Stream",
            attrs: "output"
        }        
    ]
});

blocks.register({
    name: "EntityStream",
    description: "To define an entity stream",
    fields: [
        {
            name: "SelectedType",
            type: "string",
            attrs: "editable"
        },
        {
            name: "SelectedAttributes",
            type: "string[]",
            defaultValue: ["all"],
            attrs: "editable"
        },          
        {
            name: "Groupby",
            choices: ["ALL", "EntityID", "EntityType", "EntityAttribute"],
            defaultValue: "EntityID",
            attrs: "editable"
        },
        {
            name: "Scoped",
            type: "bool",
            defaultValue: false,
            attrs: "editable"
        },         
        {
            name: "Stream",
            attrs: "output",
            type: "Stream"
        }                          
    ]
});

}



