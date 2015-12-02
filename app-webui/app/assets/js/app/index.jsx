var table = false;

$(document).ready(function () {
  var candidate = app.func.cookie('webui-sess').replace(/:/g, '":"').replace(/\|/g, ',').replace(/,/g, '","').replace(/\*/g, ':'),
      session = (candidate.length > 1) ? JSON.parse(candidate.replace('{', '{"').replace('}', '"}')) : {};
  if (session.id) {
    app.func.ajax({url: 'http://192.168.99.100:8080/authenticated', fields: {withCredentials: true}, success: function (json) {
      if ((json.header.status != 'success') || (! json.response) || (! json.response.token)) {
        return;
      }
      $('#screen-name a>img').attr('src', 'https:'+json.response.img);
    }});
  }
});

var TableRow = React.createClass({
  render: function() {
    var content = this.props.content;
    return (
        <tr key={this.props.index}>
          <td className="data-index">{content.InstanceId}</td>
        </tr>
    );
  }
});

var Table = React.createClass({
  getInitialState: function() {
    return {data: []};
  },
  load: function(sender) {
    app.func.ajax({url: '/ec2/instances/', success: function (json) {
    if (json.header.status != 'success') {
      return;
    }
    $('#count').text(json.response.count + ' instance' + ((json.response.count > 1) ? 's' : ''));
    json.response.instances.sort(function (a, b) {
      return a.id - b.id;
    });
    sender.setState({data: json.response.instances});
    }});
  },
  componentDidMount: function() {
    this.load(this);
  },
  componentWillReceiveProps: function() {
    this.load(this);
  },
  render: function() {
    var rows = this.state.data.map(function(record, index) {
      return <TableRow key={index} content={record} />
    });
    return (
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>Instance</th>
            </tr>
          </thead>
          <tbody>{rows}</tbody>
        </table>
    );
  }
});

if ($('#data').length > 0) {
  table = ReactDOM.render(<Table />, document.getElementById('data'));
}
