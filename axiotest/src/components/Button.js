import React, { Component } from 'react';

class Button extends Component {
    render() {
        return (
            <div>
                <button onClick={this.props.onDecrement}>PREV</button>
                <button onClick={this.props.onIncrement}>NEXT</button>
            </div>
        );
    }
}

export default Button;


