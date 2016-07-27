var Power = React.createClass({
    render: function() {
        return <a href="#" className="btn btn-lg btn-default"><span className="text-success glyphicon glyphicon-off" aria-label="power"></span></a>;
    }
});

ReactDOM.render(
<Power name="World" />,
    document.getElementById('container')
);