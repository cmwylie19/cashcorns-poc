// IIFE is intended to be run in the browser console on the HeyTaco leaderboard page to serve as a (rough) client application.
(() => {
  const fetchDays = () => {
    //   return fetch(`du_url_to_backend_api_to_get_days_calculated`)
    //         .then(response => {
    //                 if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);
    //                 return response.json();
    //         });

    // This is a mock of the above fetch call. It returns a promise that resolves to 15.
    console.log('Fetching days from our backend API')
    return Promise.resolve(15)
  }

  const queryForTacoCounts = days => {
    console.log(`Querying HeyTaco API with value of days as queryParameter: ${days}`)
    return fetch(
      `https://www.heytaco.chat/api/v1/json/leaderboard/T01G7FPRP8V?days=${days}`
    ).then(response => {
      if (!response.ok)
        throw new Error(`HTTP error! Status: ${response.status}`)
      return response.json()
    })
  }

  const sendTotalsToBackend = tacoCountData => {
    // return fetch(`url_to_backend_api_to_send_taco_counts`, {
    //   method: 'POST',
    //   body: JSON.stringify(tacoCountData),
    //   headers: {
    //     'Content-Type': 'application/json'
    //   }
    // })
    // .then(response => {
    //   if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);
    //   return response.json();
    // });

    // This is a mock of the above fetch call. It returns a promise that resolves to a stringified version of the taco leadboard array.
    console.log("Formatting taco count data and sending to our backend API")
    return Promise.resolve(JSON.stringify(tacoCountData.leaderboard))
  }

  const transactionSuccess = payload => {
    console.log('Taco data transfer was successful. Here is what was sent and accepted by our backend API:')
    console.log(payload)
  }

  fetchDays()
    .then(days => queryForTacoCounts(days))
    .then(tacoCountData => sendTotalsToBackend(tacoCountData))
    .then(payload => transactionSuccess(payload))
    .catch(error => console.error('There was an error!', error))

})();
