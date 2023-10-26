#!/usr/bin/env sh

#echo "start login test..."

ab -n 2 -c 1 -t 1 -T "application/json" -H "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTgxNjU0OTYsIlVzZXJJRCI6MSwiVXNlckFnZW50IjoiUG9zdG1hblJ1bnRpbWUvNy4yOS4wIn0.P7-fgPWRmO8eA_RFDV5x29kb_CyIkUkajo_7xvpFpCc" -p login_test.json http://localhost:80/users/login
