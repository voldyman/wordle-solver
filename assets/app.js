function App() {
    let self = this;
    self.query = function () {
        let req = {
            present: [
                { char: 'p', pos: 1 },
                { char: 'a', pos: 2 },
                { char: 's', pos: 5 }
            ],
            notPresent: [
                'q', 'r', 's'
            ]
        }
        fetch('/query', {
            method: 'POST',
            body: JSON.stringify(req)
        }).then((response) => {
            if (response.ok) {
                return response.json();
            }
            return Promise.reject(response);
        }).then((data) => {
            console.log(data);
        }).catch((err) => {
            console.warn(err);
        })

    }
    self.addButton = document.getElementById('addChar');
    self.addButton.onclick = function () {
        console.log('test');
        self.query();
    };


}
window.addEventListener('load', function () {
    let app = new App();
});
