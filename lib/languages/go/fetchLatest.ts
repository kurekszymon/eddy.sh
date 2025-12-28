import https from 'https';

function getLatestGoVersion() {
    return new Promise((resolve, reject) => {
        https.get('https://go.dev/dl/?mode=json', res => {
            let data = '';
            res.on('data', chunk => data += chunk);
            res.on('end', () => {
                const releases = JSON.parse(data);
                resolve(releases[0].version);
            });
        }).on('error', reject);
    });
}

getLatestGoVersion().then(console.log);