#!/bin/bash

# Get the pod name from the deployment
DEPLOYMENT_NAME="postgres"
NAMESPACE="default"

POD_NAME=$(kubectl get pods --namespace $NAMESPACE -l "app=$DEPLOYMENT_NAME" -o jsonpath="{.items[0].metadata.name}")

echo "Pod name: $POD_NAME"

# Create a SQL script with the commands
SQL_SCRIPT=$(cat <<EOF
CREATE ROLE sp SUPERUSER LOGIN PASSWORD '1234';
GRANT ALL ON SCHEMA public TO sp;
GRANT ALL ON ALL TABLES IN SCHEMA public TO sp;
CREATE DATABASE test;
\c test;
CREATE TABLE Creative_Details(
    adid varchar(20),
    height int,
    width int,
    adtype int,
    crtv_details varchar(20)
);
INSERT INTO Creative_Details VALUES
('adtest1',100,100,1,'addetails'),
('adtest2',100,50,2,'addetails'),
('adtest3',200,250,3,'crazy creative'),
('adtest4',200,150,1,'nike video add');
CREATE TABLE Budget(
    AdvID varchar(20),
    totalBudget int,
    cpm numeric(3,0),
    remBudget numeric(10,3)
);
INSERT INTO Budget VALUES
('advtest1',1000,5,1000),
('advtest2',10000,10,10000),
('advtest3',5000,6,5000),
('advtest4',4000,10,4000);
EOF
)

# Execute the SQL script inside the pod
kubectl exec -i $POD_NAME -- psql -U postgres <<EOF
$SQL_SCRIPT
EOF

echo "PostgreSQL setup completed."