var ws = new WebSocket('ws://localhost:8080/ws');

class Power extends React.Component {
    constructor() {
        super();
        this.handleClick = this.handleClick.bind(this);
        this.state = {
            "status": "ON"
        }
    }

    componentDidMount() {
        var component = this;
        ws.onmessage = function (evt)
        {
            var received_msg = evt.data;
            component.setState({status: received_msg});
            console.log(received_msg);
        };
    }

    handleClick() {
        var newStatus = this.state.status == "ON" ? "OFF" : "ON";
        fetch("http://localhost:8080/power", {
            method: "PUT",
            body: JSON.stringify({
                state: newStatus
            })
        })
    }

    render() {
        var className = (this.state.status == "ON" ? "text-success" : "text-danger");
        className = className + " glyphicon glyphicon-off";
        return (
            <a href="#" className="btn btn-lg btn-default" onClick={this.handleClick}>
                <span className={className} aria-label="power"></span>
            </a>
        );
    }
}

ReactDOM.render(
<Power initialStatus="ON" />,
    document.getElementById('container')
);