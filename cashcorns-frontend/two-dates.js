const DaysSinceDate = (inputDate) => {

    const inputDateObj = new Date(inputDate);

    if (isNaN(inputDateObj)) {
        return "Invalid Date";
    }

    const currentDate = new Date();
    const timeDifference = currentDate - inputDateObj;

    const daysDifference = Math.floor(timeDifference / (1000 * 60 * 60 * 24));

    return daysDifference;
}

const fetchTacos = async (days) => {
    let users = {};
    return fetch(`https://www.heytaco.chat/api/v1/json/leaderboard/T01G7FPRP8V?days=${days}`)
        .then(res => res.json())
        .then(data => {
            // change sum to dollar amount
            const TACO_VALUE = 5;
            data.leaderboard.map((user) => {

                user.sum = `${user.sum * TACO_VALUE}`
                users[user.username] = user
            })
            return users;
        })
}


(async () => {
    // DATE_1 is the earliest date
    const DATE_1 = "2023-09-27";

    // DATE_2 is the most recent date
    const DATE_2 = "2023-09-28";

    const DAYS_SINCE_DATE_1 = DaysSinceDate(DATE_1);
    const DAYS_SINCE_DATE_2 = DaysSinceDate(DATE_2);

    console.log(`Days since ${DATE_1}: ${DAYS_SINCE_DATE_1}`)
    console.log(`Days since ${DATE_2}: ${DAYS_SINCE_DATE_2}`)

    const LEADERBOARD_DATE_1 = await fetchTacos(DAYS_SINCE_DATE_1);
    const LEADERBOARD_DATE_2 = await fetchTacos(DAYS_SINCE_DATE_2);

    let final_array = {};
    final_array.leaderboard = [];
    // Remove the sum from the most recent date and today to find window
    for (const username in LEADERBOARD_DATE_1) {
        // only for people in both dates
        if (LEADERBOARD_DATE_1.hasOwnProperty(username) && LEADERBOARD_DATE_2.hasOwnProperty(username)) {
            LEADERBOARD_DATE_1[username].sum = `${LEADERBOARD_DATE_1[username].sum - LEADERBOARD_DATE_2[username].sum}`
        }
        final_array.leaderboard.push(LEADERBOARD_DATE_1[username])
    }



    // Create a Blob containing the JSON data
    const blob = new Blob([JSON.stringify(final_array, undefined, 2)], { type: 'application/json' });

    // Create a temporary URL for the Blob
    const blobUrl = window.URL.createObjectURL(blob);

    // Create an anchor element for the download
    const a = document.createElement('a');
    a.style.display = 'none';
    a.href = blobUrl;
    a.download = `Taco-Window-${DATE_1}-${DATE_2}.json`;

    // Append the anchor element to the document and trigger the click event
    document.body.appendChild(a);
    a.click();

    // Clean up: remove the anchor and revoke the Blob URL
    document.body.removeChild(a);
    window.URL.revokeObjectURL(blobUrl);


})()
