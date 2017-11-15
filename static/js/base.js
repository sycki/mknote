

new Vue({
	el: '#sections',
    data:{
        items:[
            {
                name:'家用电器',
                isSubShow:false,
                subItems:[
                    {
                        name:'笔记本电脑'
                    },
                    {
                        name:'台式电脑'
                    },
                    {
                        name:'电视机'
                    }
                ]
            },
            {
                name:'服装',
                isSubShow:false,
                subItems:[
                    {
                        name:'男士服装'
                    },
                    {
                        name:'女士服装'
                    },
                    {
                        name:'青年服装'
                    }
                ]
            }
        ]
    },
    methods:{
        showToggle:function(item){
            item.isSubShow = !item.isSubShow
        }
    }
})  