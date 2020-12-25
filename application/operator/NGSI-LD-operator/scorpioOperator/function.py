ctxELement = {}

# handle notify Entity
def handleEntity(ctxObj, create, update, append):
    print('===============Implement losic====================')
    print(ctxObj)
    for ctx in ctxObj:
	handleScorpioUpdate(ctx, create, update, append)


#handle update and append request for scorpio broker

def handleupdateAppend(currUpdateCtx, create, update, append):    
    appendCtx = {}
    global ctxELement
    eid = currUpdateCtx['id']
    preCtxEle = ctxELement[eid]
    appendCtx['id'] = currUpdateCtx['id']
    appendCtx['type'] = currUpdateCtx['type']
    for key in  currUpdateCtx : 
        if ctxELement[eid].has_key(key) == False :
	    appendCtx[key] = currUpdateCtx[key]
	    preCtxEle[key] = currUpdateCtx[key]
    
    ctxELement[eid] = preCtxEle
    if len(appendCtx) > 0: 
        append(appendCtx)
    update(currUpdateCtx)
    

# handle creation of etity on scorpio broker

def handleScorpioUpdate(ctx, create, update, append):
	global ctxELement 
	eid = ctx['id']
	if ctxELement.has_key(eid) == True:
	    handleupdateAppend(ctx, create, update, append)
	else:
	    ctxELement[eid] = ctx
	    create(ctx)


def handleAlreadyCreatedEntity(eid, create, update, append):
	global ctxELement
	ctxObj = ctxELement[eid]
	print("ctxObj")
        print(ctxObj)
	handleupdateAppend(ctxObj,create, update, append)
        

