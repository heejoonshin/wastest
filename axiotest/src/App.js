import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';

import Button from './components/Button';
import * as service from './services/todolist';

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            id: 0,

        };
    }

    onIncrement = (event) => {
        this.fetchUserInfo(this.state.id + 1);
    }

    onDecrement = (event) => {
        this.fetchUserInfo(this.state.id - 1);
    }

    fetchUserInfo = async (id) => {
        const info = await Promise.all([
            service.GetTodo(id),
            service.GetTodolist({limit:3})


        ]);

        this.setState(prevState => ({
            id: id,
            title: info.keys()

        }));
    }

    render() {
        return (
            <div>
                <h1>{this.state.title}</h1>
                <p>{this.state.title}</p>
                <Button onIncrement={this.onIncrement} onDecrement={this.onDecrement} />
            </div>
        );
    }
}

export default App;

