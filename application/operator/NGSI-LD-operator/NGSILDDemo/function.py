def handleEntity(ctxObj, publish):
    print('===============Implement losic====================')
    print(ctxObj)
    for ctx in ctxObj :
        if ctx.has_key('temprature'):
            temprature = ctx['temprature']
	    print("Temprature")
            print(temprature)
            if temprature.has_key('value') and (temprature['value'] > 50):
		    publish(ctx)
    
    
