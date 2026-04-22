name: Deploy Backend

on:
  push:
    branches:
      - main
      - develop

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Set Environment
        id: env
        run: |
          if [[ "${GITHUB_REF}" == "refs/heads/main" ]]; then
            echo "ENV=production" >> $GITHUB_OUTPUT
            echo "TAG=latest" >> $GITHUB_OUTPUT
          else
            echo "ENV=staging" >> $GITHUB_OUTPUT
            echo "TAG=staging" >> $GITHUB_OUTPUT
          fi

      - name: Login to GHCR
        run: echo "${{ secrets.GITHUB_TOKEN }}" | docker login ghcr.io -u ${{ github.actor }} --password-stdin

      - name: Build Docker Image
        run: |
          docker build -t ghcr.io/${{ github.repository }}:${{ steps.env.outputs.TAG }} .

      - name: Push Docker Image
        run: |
          docker push ghcr.io/${{ github.repository }}:${{ steps.env.outputs.TAG }}

      - name: Deploy via SSH
        uses: appleboy/ssh-action@v1.2.0
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SERVER_SSH_KEY }}
          script: |
            docker pull ghcr.io/${{ github.repository }}:${{ steps.env.outputs.TAG }}

            docker stop go-backend || true
            docker rm go-backend || true

            docker run -d \
              --name go-backend \
              -p 8080:8080 \
              --restart always \
              ghcr.io/${{ github.repository }}:${{ steps.env.outputs.TAG }}