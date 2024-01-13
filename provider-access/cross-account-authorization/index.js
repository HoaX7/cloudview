import url from "url"

const STATUS = {
  SUCCESS: "SUCCESS",
  FAILED: "FAILED"
}

export const handler = async (event, context, callback) => {
  console.log("Event Data:", event)
  await sendResponse(event, context, STATUS.SUCCESS, {})
  return
};

// Send response to the pre-signed S3 URL
async function sendResponse(event, context, responseStatus, responseData = {}) {
  var responseBody = JSON.stringify({
    Status: responseStatus,
    Reason:
      "See the details in CloudWatch Log Stream: " + context.logStreamName,
    PhysicalResourceId: context.logStreamName,
    StackId: event.StackId,
    RequestId: event.RequestId,
    LogicalResourceId: event.LogicalResourceId,
    Data: responseData,
  });

  console.log("RESPONSE BODY:\n", responseBody);

  await putCFResponse(event.ResponseURL, responseBody, context);
  return
}

const putCFResponse = async (responseURL, responseBody, context) => {
  var parsedUrl = url.parse(responseURL);
  var options = {
    hostname: parsedUrl.hostname,
    port: 443,
    pathname: parsedUrl.path,
    method: "PUT",
    body: responseBody
  };

  console.log("SENDING RESPONSE to CF stack...\n");
  try {
    const response = await fetch(options, {
      headers: {
        "content-type": "",
        "content-length": responseBody.length,
      },
    })
    console.log("Status:", response.status)
    console.log("Headers:", JSON.stringify(response.headers))
  } catch (err) {
    console.log("Fetch request failed..", err)
  }
  return
  // var request = https.request(options, function (response) {
  //   console.log("STATUS: " + response.statusCode);
  //   console.log("HEADERS: " + JSON.stringify(response.headers));
  //   // Tell AWS Lambda that the function execution is done
  //   // NO NEED TO USE IF YOU WANT TO USE ASYNC/AWAIT
  //   context.done();
  // });

  // request.on("error", function (error) {
  //   console.log("sendResponse Error:" + error);
  //   // Tell AWS Lambda that the function execution is done
  //   // NO NEED TO USE IF YOU WANT TO USE ASYNC/AWAIT
  //   context.done();
  // });

  // write data to request body
  // request.write(responseBody);
  // request.end();
}
