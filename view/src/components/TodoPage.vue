<template>
  <div class="container">
    <h2>Todo List</h2>


    <b-table :items="todolist" :fields="fields" >
      <template slot="modify" slot-scope="row">
        <b-button size="lg" @click.stop="row.toggleDetails" class="mr-2">
          수정
        </b-button>

      </template>
      <template slot="delete" slot-scope="row">
        <b-button size="lg"  class="mr-3" type="button" @click="deleteTodo(row.item.Id)" @click.stop="reload_data">
          삭제
        </b-button>

      </template>

      <template slot="row-details" slot-scope="row">
        <b-card>
          <b-row class="mb-2">
            <b-col sm="3" class="text-sm-right">
            <input type="text" class="form-control"
                   placeholder="할일을 입력하세요"
                   v-model="M_title"></b-col>
          </b-row>
          <b-row class="mb-2">
            <b-col sm="3" class="text-sm-right">
            <input type="text" class="form-control"
                   placeholder="참조할 일을 입력하세요 ex: 1 2 3 1번과 2번과 3번이 참조"
                   v-model="M_ref"></b-col>
          </b-row>
          <b-row class="mb-2">
            <b-col sm="3" class="text-sm-right">
              <b-form-group label="완료처리">
                <b-form-radio-group id="btnradios2"
                                    buttons
                                    button-variant="outline-primary"
                                    size="lg"
                                    v-model="selected"
                                    :options="options"
                                    name="radioBtnOutline" />
              </b-form-group>

            </b-col>


          </b-row>
          <b-button size="sm" @click="reload_data" @click.stop="modifyTodo({title:M_title,children:M_ref,done:selected},row.item.Id)"  >수정 완료</b-button>

        </b-card>
      </template>

    </b-table>
    <div id="example">{{ message }}</div>
    <b-pagination size="md" :total-rows="totaldata" :per-page="limit" :limit="10" v-model="currentPage" @input="getTodolist({page:currentPage})"></b-pagination>
    <div class="input-group" style="margin-bottom:10px;">
      <input type="text" class="form-control"
             placeholder="할일을 입력하세요"
             v-model="name">
      <input type="text" class="form-control"
             placeholder="참조할 일을 입력하세요 ex: 1 2 3 1번과 2번과 3번이 참조"
             v-model="ref">
      <span class="input-group-btn">
		<button class="btn btn-default" type="button"
            @click="createTodo({title:name,children:ref})" @click.stop="reload_data">추가</button>
	</span>
    </div>
  </div>
</template>


<script>


  export default {


    name: 'TodoPage',
    data() {
      return {
        todolist : [],
        totaldata : '',
        limit : '',
        fields: [
          {
            key: 'Id',
            sortable: true
          },
          {
            key: 'Title',
            sortable: true,
            label: "할일"

          },
          {
            key: 'CreatedAt',
            sortable: true,
            label: "작성일시"
          },
          {
            key: 'UpdatedAt',
            label: '최종수정일시',
            sortable: true,

          },
          {
            label: '완료처리',
            key: 'Done',
          },
          {
            label:"수정",
            key:'modify'
          },
          {
            label:"삭제",
            key:'delete'
          }
        ],
        options: [
          { text: '완료', value: 'Y' },
          { text: '미완료', value: 'N' },

        ],
        page : 1
      }

    },
    created:function () {
      this.$nextTick(function (){

      })
    },

    methods: {




      deleteTodo(id) {

        if(id != null){
          var vm = this;
          this.$http.defaults.headers.post['Content-Type'] = 'application/json';
          this.$http.delete('http://localhost:8080/api/todo/'+id).then((result) => {
            console.log(result);
            for (var i = 0; i < this.todolist.length; i++){


              if(this.todolist[i].Id == result.data.Id)
              {
                this.todolist.splice(i,1)

              }
            }
          }).catch(e=>{
            alert(e.response.data.Message)

          });
          this.name = null
        }
      },
      getTodolist(params) {



          const baseURI = 'http://localhost:8080';
          this.$http.get(baseURI + "/api/todolist/", {
            params
          })
            .then((result) => {
              this.todolist = result.data.Todolist;
              this.totaldata = result.data.Totaldata;
              this.limit = result.data.Limit
              this.page = result.data.Page


            }).catch((error) => {
            alert(error.response.data.Message)

          })

      },
      createTodo(params){
        console.log(params)


        if(params.title != null){
          var vm = this;
          if(params.children != null){
            params.children = params.children.trim();

            params.children = params.children.split(' ')

            if(params.children != "") {
              for (var i = 0; i < params.children.length; i++) {
                params.children[i] = parseInt(params.children[i])
                console.log(typeof params.children[i]);


              }
            }else{
              params.children=null

            }

          }
          console.log(params);
          this.$http.defaults.headers.post['Content-Type'] = 'application/json';
          this.$http.post('http://localhost:8080/api/todo/',params).then((result) => {
            console.log(result);
            vm.todolist.push(result.data)
          }).catch((e)=>{
            alert(e.response.data.Message)


          });
          this.name = null
        }
      },
      reload_data(){
        window.location.reload()
      },
      modifyTodo(params,id){
        console.log(params)
        if(params != null){
          if(params.title != null){
            let x = params.title.trim()
            if(x == ''){
              params.title = null
            }
          }
          if(params.children != null) {
            params.children = params.children.trim();

            params.children = params.children.split(' ')

            if (params.children != "") {
              for (var i = 0; i < params.children.length; i++) {
                params.children[i] = parseInt(params.children[i])
                console.log(typeof params.children[i]);


              }
            } else {
              params.children = null

            }
          }

          var vm = this;
          this.$http.defaults.headers.post['Content-Type'] = 'application/json';
          this.$http.put('http://localhost:8080/api/todo/'+id,params).then((result) => {
            console.log(result);
            this.$router.go({
              path: $router.path,

            });


            for (var i = 0; i < this.todolist.length; i++){

              if(this.todolist[i].Id == result.data.Id) {
                vm.todolist[i] = result.data;

                console.log(vm.todolist[i]);


                break;
              }

            }



          }).catch(e=>{

            alert(e.response.data.Message)

          })
        }
      }

    },
    mounted(){
      this.getTodolist()
    }
  }
</script>
