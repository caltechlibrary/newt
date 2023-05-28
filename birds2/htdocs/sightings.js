/* sightings.js provides access to our JSON API run by PostgREST
   and assembles the results before updating the web page. */
(function(document, window) {
  let list_url = 'http://localhost:3000/bird_view',
    record_url = 'http://localhost:3000/rpc/record_bird',
    list_elem = document.getElementById('bird-list'),
    bird_elem = document.getElementById('bird'),
    place_elem = document.getElementById('place'),
    sighted_elem = document.getElementById('sighted'),
    add_button = document.getElementById('record-bird');

  function updateList(elem, src) {
    let bird_list = JSON.parse(src),
      parts = [];
    parts.push('<table>');
    parts.push('  <tr><th>bird</th><th>place</th><th>sighted</th></tr>');
    for (const obj of bird_list) {
      parts.push(` <tr><td>${obj.bird}</td><td>${obj.place}</td><td>${obj.sighted}</td></td>`);
    }
    parts.push('</table>');
    elem.innerHTML = parts.join('\n');
  }

  function birdRecord(bird_elem, place_elem, sighted_elem) {
    return { "bird": bird_elem.value, "place": place_elem.value, "sighted": sighted_elem.value };
  }

  function getData(elem, url, updateFn) {
    /* We use a xhr to retrieve the current list of sightings. */
    const req = new XMLHttpRequest();
    req.addEventListener("load", function(evt) {
      /* Call our page update function */
      updateFn(elem, this.responseText);
    });
    req.open("GET", url);
	req.setRequestHeader('Cache-Control', 'no-cache');
    req.send();
  };

  function postData(obj, url) {
    const req = new XMLHttpRequest();
    req.open("POST", url, true);
    req.setRequestHeader('Content-Type', 'application/json');
    req.onreadystatechange = function() {//Call a function when the state changes.
      console.log(`DEBUG state ${req.readyState}, status ${req.status}, resp ${req.responseText}`);
    }
    req.send(JSON.stringify(obj));
  }

  /* Main processing for page */
  add_button.addEventListener("click", function(evt) {
    postData(birdRecord(bird_elem, place_elem, sighted_elem), record_url);
	/* Now we need to update our listing! */
	list_elem.innerHTML = '';
	setTimeout(() => {
  		console.log("Delayed for 10 second.");
  		getData(list_elem, list_url, updateList);
	}, "10 second");
    evt.preventDefault();
  });

  getData(list_elem, list_url, updateList);
})(document, window);
