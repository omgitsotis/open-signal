$(document).ready(function(){
    $("#tower-form").submit(function(event){
        event.preventDefault();
        let lng = $("input[name=lng]").val();
        let lat = $("input[name=lat]").val();

        let $tableBody = $('#results-body');
        $tableBody.empty();

        let url = "/towers?lat=" + lat + "&lng=" + lng;
        fetch(url)
            .then(response => response.json())
            .then(towerList => {
                towerList.forEach(tower => {
                    row = document.createElement("tr");

                    Object.keys(tower).forEach(function(key) {
                        let td = document.createElement("td");
                        td.innerHTML = tower[key]
                        row.appendChild(td);
                    });

                    $tableBody.append(row);
                })
            });
    });
});
