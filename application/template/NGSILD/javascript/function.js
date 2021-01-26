//
//  contextEntity: the received entities
//  publish, query, and subscribe are the callback functions for your own function to interact with the assigned nearby broker
//      publish:    publish the generated entity, which could be a new entity or the update of an existing entity
//      query:      query addtional information from the assigned nearby broker
//      subscribe:  subscribe addtional infromation from the assigned nearby broker
//
exports.handler = function(contextEntity, publish, query, subscribe)
{
	console.log("enter into the user-defined fog function");
	console.log(contextEntity)	
	if (contextEntity == null) {
		return;
	} 	
	
	publish(contextEntity)
};

