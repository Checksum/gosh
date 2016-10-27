package handlers

const ScriptContent string = `
<script id="_gosh_reloader">
(function() {
    var source = new EventSource("/__events");
    var retries = 0;
    
    source.onopen = function(e) {
        retries = 0;
    };

    source.onmessage = function(e) {
        window.location.reload(true);
    };

    source.onerror = function() {
        retries += 1;
        if (retries >= 3) {
            source.close();
        }        
    };
}());
</script>
`
