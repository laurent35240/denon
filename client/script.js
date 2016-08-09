class Power extends React.Component {
    constructor() {
        super();
        this.handleClick = this.handleClick.bind(this);
    }
    handleClick() {
        fetch("http://localhost:8080/power", {
            method: "PUT",
            body: {
                state: "ON"
            }
        })
    }
    render() {
        return (
            <a href="#" className="btn btn-lg btn-default" onClick={this.handleClick}>
                <span className="text-success glyphicon glyphicon-off" aria-label="power"></span>
            </a>
        );
    }
}

ReactDOM.render(
<Power name="World" />,
    document.getElementById('container')
);