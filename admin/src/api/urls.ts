export const getAbsoluteUrl = (relativeUrl: string): string => {
    const hostname = window.location.hostname;
    return `http://${hostname}:9401${relativeUrl}`;
};
