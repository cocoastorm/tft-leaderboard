FROM gcr.io/distroless/base-debian11
ADD ./site/out /www
ADD ./tft-leaderboard /leaderboard

# run as web server
CMD ["/leaderboard", "serve", "--app-path", "/www"]
