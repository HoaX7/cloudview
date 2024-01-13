const url = "https://cloudformation-custom-resource-response-useast1.s3.amazonaws.com/arn%3Aaws%3Acloudformation%3Aus-east-1%3A009771722451%3Astack/cloudfriendly-access-clr6ep9yu002kybpw4n6n4j0r/364e8210-af85-11ee-b6a4-0a1c4fc4028d%7CCloudFriendlyAuthenticator%7Cb71f4eec-198e-46da-a3a1-3b04e8fa2f6a?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20240110T065541Z&X-Amz-SignedHeaders=host&X-Amz-Expires=7200&X-Amz-Credential=AKIA6L7Q4OWTWMBLHBGG%2F20240110%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Signature=19844980ad3f04ad0acc0089e2bdbf1b6f46aa983686fcea324c0d348ce6b811"

import URL from "url"
// const https = require("https")
import fetch from "node-fetch"
const boot = async () => {
    var responseBody = JSON.stringify({
      "Status": "SUCCESS",
      "Reason": "See the details in CloudWatch Log Stream: 2024/01/10/[$LATEST]1af047c0541342688e8f380a836024a9",
      "PhysicalResourceId": "2024/01/10/[$LATEST]1af047c0541342688e8f380a836024a9",
      "StackId": "arn:aws:cloudformation:us-east-1:009771722451:stack/cloudfriendly-access-clr6ep9yu002kybpw4n6n4j0r/364e8210-af85-11ee-b6a4-0a1c4fc4028d",
      "RequestId": "b71f4eec-198e-46da-a3a1-3b04e8fa2f6a",
      "LogicalResourceId": "CloudFriendlyAuthenticator",
      "Data": {}
  });
  var options = {
    port: 443,
    method: "PUT",
    body: responseBody,
    headers: {
      "content-type": "",
      "content-length": responseBody.length,
    },
  };

  console.log("SENDING RESPONSE...\n");
  console.log(options, url)
  /*global fetch*/
  await fetch(url, options)
  // var request = https.request(options, function (response) {
  //   console.log("STATUS: " + response.statusCode);
  //   console.log("HEADERS: " + JSON.stringify(response.headers));
  //   // Tell AWS Lambda that the function execution is done
  // });

  // request.on("error", function (error) {
  //   console.log("sendResponse Error:" + error);
  //   // Tell AWS Lambda that the function execution is done
  // });

  // // write data to request body
  // request.write(responseBody);
  // request.end();
}

boot()
