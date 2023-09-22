function CalculateDaysSinceMostRecentDate(dateArray) {

    const today = new Date();

    // only need past dates
    const pastDates = dateArray.filter(dateString => {
        const date = new Date(dateString);
        return date < today;
    });


    if (pastDates.length === 0) {
        return null;
    }

    // sort
    pastDates.sort((a, b) => new Date(b) - new Date(a));

   
    const mostRecentDate = new Date(pastDates[0]);
    const timeDifference = today - mostRecentDate;
    const daysSinceMostRecent = Math.floor(timeDifference / (1000 * 60 * 60 * 24));

    return daysSinceMostRecent;
}

(() => {
    const approvalDates = [
        "2023-09-19",
        "2023-09-27",
        "2023-10-11"
    ];

    days = CalculateDaysSinceMostRecentDate(approvalDates)
    url = `https://www.heytaco.chat/api/v1/json/leaderboard/T01G7FPRP8V?days=${days}`
    fetch(url).then(res => res.json()).then(data => {
        console.log(data, undefined, 2)
        return data;
    }).then(data => {

        // Create a Blob containing the JSON data
        const blob = new Blob([JSON.stringify(data,undefined,2)], { type: 'application/json' });

        // Create a temporary URL for the Blob
        const blobUrl = window.URL.createObjectURL(blob);

        // Create an anchor element for the download
        const a = document.createElement('a');
        a.style.display = 'none';
        a.href = blobUrl;
        a.download = `Taco-${new Date().toDateString().replaceAll(" ","-")}.json`;

        // Append the anchor element to the document and trigger the click event
        document.body.appendChild(a);
        a.click();

        // Clean up: remove the anchor and revoke the Blob URL
        document.body.removeChild(a);
        window.URL.revokeObjectURL(blobUrl);
    })

})()
