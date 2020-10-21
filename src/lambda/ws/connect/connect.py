import json

import boto3

def handle(event, context):
    connections_table = boto3.resource('dynamodb').Table('connections')
    connection_id = event['requestContext']['connectionId']
    connections_table.put_item(Item={
        'ConnectionID': connection_id,
        'timestamp': event['requestContext']['connectedAt']
        })
    domain_name = event['requestContext']['domainName']
    stage = event['requestContext']['stage']
    endpoint = f'https://{domain_name}'
    agw_session = boto3.session.Session()
    agw_api = agw_session.client(service_name='apigatewaymanagementapi', endpoint_url=endpoint)
    client_payload = {
        'poolId': '123'
    }
    try:
        agw_api.post_to_connection(
            Data=json.dumps(client_payload).encode(),
            ConnectionId=event['requestContext']['connectionId']
        )
    except Exception as e:
        print(str(e))
    return {
        'statusCode': 200
    }
