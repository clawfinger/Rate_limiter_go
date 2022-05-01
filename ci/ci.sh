make lint
if [ $? -eq 0 ] 
then 
  echo "Lint success" 
else 
  echo "Lint failed"
  exit 1 
fi
make test
if [ $? -eq 0 ] 
then 
  echo "Tests success" 
else 
  echo "Tests failed"
  exit 1 
fi
make build
if [ $? -eq 0 ] 
then 
  echo "Build success" 
else 
  echo "Build failed"
  exit 1 
fi
