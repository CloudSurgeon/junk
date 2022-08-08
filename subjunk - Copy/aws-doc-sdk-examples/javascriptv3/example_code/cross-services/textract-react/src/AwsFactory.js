// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

import { CognitoIdentityClient } from "@aws-sdk/client-cognito-identity";
import { fromCognitoIdentityPool } from "@aws-sdk/credential-provider-cognito-identity";
import { TextractClient } from "@aws-sdk/client-textract";
import { S3Client } from "@aws-sdk/client-s3";
import { SQSClient } from "@aws-sdk/client-sqs";

/**
 * A factory that creates AWS clients used by the demo application.
 * This application requires an authenticated Amazon Cognito identity to grant
 * permissions.
 *
 * @param cognitoId: The ID of the Amazon Cognito identity provider for the user pool.
 * @param cognitoToken: The token for the session generated by Amazon Cognito.
 * @param cognitoIdentityPoolId: The ID of the Amazon Cognito identity pool that
 *                               contains the logged in user.
 * @param deployRegion: The AWS Region where the AWS resources are deployed for the
 *                      demo.
 * @returns {{}}: The AWS clients, initialized with Amazon Cognito identity credentials.
 */
export const awsFactory = ({
  cognitoId,
  cognitoToken,
  cognitoIdentityPoolId,
  deployRegion,
}) => {
  const creds = fromCognitoIdentityPool({
    client: new CognitoIdentityClient({ region: deployRegion }),
    identityPoolId: cognitoIdentityPoolId,
    logins: { [cognitoId]: cognitoToken },
  });

  const clients = {};
  clients["s3"] = new S3Client({
    region: deployRegion,
    credentials: creds,
  });

  clients["textract"] = new TextractClient({
    region: deployRegion,
    credentials: creds,
  });

  clients["sqs"] = new SQSClient({
    region: deployRegion,
    credentials: creds,
  });

  return clients;
};
