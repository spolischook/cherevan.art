import axios from 'axios';
import AWS from 'aws-sdk';

const s3 = new AWS.S3();

export const handler = async (event) => {
    const url = 'https://cherevan.art';

    try {
        const response = await axios.get(url);
        const data = response.data;

        // Upload data to S3
        const params = {
            Bucket: 'test-internet-and-aws-resources-access',
            Key: 'cherevan.art.html',
            Body: data
        };

        await s3.upload(params).promise();

        return {
            statusCode: 200,
            body: 'Data uploaded to S3 successfully'
        };
    } catch (error) {
        console.error('Error:', error);
        return {
            statusCode: 500,
            body: 'Error retrieving or uploading data: ' + error.message
        };
    }
};
