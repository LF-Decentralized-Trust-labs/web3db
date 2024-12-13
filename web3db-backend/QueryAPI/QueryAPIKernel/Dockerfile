# Base image for the application runtime
FROM mcr.microsoft.com/dotnet/aspnet:6.0 AS base
WORKDIR /app
# Expose HTTP and HTTPS ports
EXPOSE 80
EXPOSE 443

# Image used to build the application
FROM mcr.microsoft.com/dotnet/sdk:6.0 AS build
WORKDIR /src
# Copy the project file and restore dependencies
COPY ["1.1-Presentation/Query.API.Kernel/Query.API.Kernel.csproj", "1.1-Presentation/Query.API.Kernel/"]
RUN dotnet restore "1.1-Presentation/Query.API.Kernel/Query.API.Kernel.csproj"
# Copy the remaining source code and build the project
COPY . .
WORKDIR "/src/1.1-Presentation/Query.API.Kernel"
RUN dotnet build "Query.API.Kernel.csproj" -c Release -o /app/build

# Publish the application
FROM build AS publish
RUN dotnet publish "Query.API.Kernel.csproj" -c Release -o /app/publish /p:UseAppHost=false

# Final image for running the application
FROM base AS final
WORKDIR /app
# Copy the published output from the previous stage
COPY --from=publish /app/publish .

# Create a directory for the certificates
RUN mkdir -p /https

# Environment variables to configure Kestrel for HTTPS
ENV ASPNETCORE_URLS="http://+:80;https://+:443"
ENV ASPNETCORE_Kestrel__Certificates__Default__Path=/https/fullchain.pem
ENV ASPNETCORE_Kestrel__Certificates__Default__KeyPath=/https/privkey.pem

# Set the entry point for the application
ENTRYPOINT ["dotnet", "Query.API.Kernel.dll"]