import axios from 'axios';

export function GetTodo(id){
    return axios.get('http://localhost:8080/api/todo/' + id)
}

export function GetTodolist(params){
    return axios.get('http://localhost:8080/api/todolist' )
}
