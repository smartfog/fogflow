(function () {

    /**
     * check data exists in list
     * @param {string} data : variable name
     * @param {list} list_ : list of response data
     * @returns list of filtered data
     */
    function isDataExists(data, list_) {

        var result = list_.filter(x => x.name === data);
        console.log("is data exists ",result);
        return result;
    }



    if (typeof module !== 'undefined' && typeof module.exports !== 'undefined') {
        this.axios = require('axios')
        module.exports.isDataExists = isDataExists;
        module.exports.testfunction = testfunction;
    } else {
        window.isDataExists = isDataExists;
        window.testfunction = testfunction;
    }
})();